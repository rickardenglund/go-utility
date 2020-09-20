package workpool_test

import (
	"fmt"
	"go-utility/workpool"
	"testing"

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
