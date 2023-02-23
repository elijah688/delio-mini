package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/constraints"
)

type Suite struct {
	suite.Suite
}

func Max[T constraints.Ordered](a T, bs ...T) T {
	max := a
	for _, b := range bs {
		if b > a {
			max = b
		}
	}

	return max
}
func TestRun(t *testing.T) {
	suite.Run(t, new(Suite))
}

func BenchmarkIntMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Max(i, []int{1, 2, 3, 4, 5, 6}...)
		// main()
	}
}
