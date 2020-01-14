package itertools

import (
	"reflect"
	"testing"
)

type point struct {
	x int
	y int
}

func TestStructIteration(t *testing.T) {
	m := []point{
		point{1, 2},
		point{3, 4},
		point{2, 4},
	}
	expected := []point{
		point{1, 2},
	}
	points := []point{}
	SliceIter(m).
		Filter(func(p point) bool {
			return p.x == 1
		}).
		Each(func(p point) {
			points = append(points, p)
		})
	if !reflect.DeepEqual(points, expected) {
		t.Error("filtered slice contains not expected result")
	}
}
