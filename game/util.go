package game

import "math"

func distanceBetween(one hasPosition, two hasPosition) float64 {

	onex, oney := one.GetPosition()
	twox, twoy := two.GetPosition()

	first := math.Pow(float64(twox-onex), 2)
	second := math.Pow(float64(twoy-oney), 2)

	return math.Sqrt(first + second)
}
