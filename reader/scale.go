package reader

// wow beautiful code <3
func scaleDistance(lastAverage, newAverage float64) int {
	difference := newAverage - lastAverage

	switch {
	case difference <= 500 && difference >= -500:
		return 3
	case difference > 500 && difference <= 1000:
		return 4
	case difference > 1000:
		return 5
	case difference < -500 && difference >= -1000:
		return 2
	case difference < -1000:
		return 1
	default:
		return 3
	}
}
