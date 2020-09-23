package workpool_test

import (
	"fmt"
	"testing"

	"github.com/rickardenglund/go-utility/workpool"

	"github.com/stretchr/testify/require"
)

func TestDoWork(t *testing.T) {
	wp := workpool.New(5)
	work := []int{1, 2, 3, 4, 5}

	wp.DoParallel(len(work),
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

	wp.DoParallel(len(work),
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

func TestName(t *testing.T) {
	res := make([]int, 3)
	wp := workpool.New(1)
	wp.DoParallel(3, func(workIndex int) {
		res[workIndex] = workIndex
	})

	require.Equal(t, []int{0, 1, 2}, res)

}
