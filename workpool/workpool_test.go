package workpool_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rickardenglund/go-utility/workpool"

	"github.com/stretchr/testify/require"
)

func TestDoWork(t *testing.T) {
	wp := workpool.New(5)
	work := []int{1, 2, 3, 4, 5}

	wp.DoParallel(context.Background(), len(work),
		func(i int) {
			work[i] += 1
		},
	)

	require.Equal(t, []int{2, 3, 4, 5, 6}, work)
}

func TestDoComplex(t *testing.T) {
	type Response struct {
		r   string
		err error
	}

	// ACT
	wp := workpool.New(5)

	work := []string{"a", "aa", "b"}
	res := make([]Response, len(work))

	wp.DoParallel(context.Background(), len(work),
		func(i int) {
			res[i].r, res[i].err = doer(work[i])
		},
	)

	// ASSERT
	require.Equal(t, "aa", res[0].r)
	require.Errorf(t, res[1].err, "should be error since len is 2")
}

func doer(s string) (string, error) {
	if len(s) == 2 {
		return "", fmt.Errorf("len is 2")
	}

	return s + s, nil
}

func TestAllworkIsDone(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res := make([]int, 3)
	wp := workpool.New(1)

	wp.DoParallel(ctx, 3, func(workIndex int) {
		res[workIndex] = workIndex
	})

	require.Equal(t, []int{0, 1, 2}, res)
}

func TestCancelTerminates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	res := make([]int, 3)
	wp := workpool.New(1)

	wp.DoParallel(ctx, 3, func(workIndex int) {
		fmt.Printf("sleeping: %v\n", time.Now().Nanosecond())
		time.Sleep(time.Duration((workIndex+1)*50) * time.Millisecond)
		res[workIndex] = workIndex + 1
	})

	require.NotEqual(t, []int{1, 2, 3}, res)
}
