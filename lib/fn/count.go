package fn

func CountMapValues[K comparable, V comparable](m map[K]V) map[V]int {
	counts := make(map[V]int)
	for _, v := range m {
		counts[v]++
	}
	return counts
}