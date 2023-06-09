package pointsconvertor

func FromPenny(points int) float64 {
	return float64(points) / 100
}

func ToPenny(points float64) int {
	return int(points * 100)
}
