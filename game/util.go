package game

import (
	"math"
)

// distanceBetween calculates distance using Chebyshev distance
// 	https://en.wikipedia.org/wiki/Chebyshev_distance
func distanceBetween(one hasPosition, two hasPosition) float64 {

	x1, y1 := one.GetPosition()
	x2, y2 := two.GetPosition()

	x := math.Abs(float64(x2 - x1))
	y := math.Abs(float64(y2 - y1))

	return math.Max(x, y)
}

// 4 [ ][ ][ ][c]
// 3 [ ][ ][b][ ]
// 2 [ ][a][ ][ ]
// 1 [ ][ ][ ][ ]
//    1  2  3  4
