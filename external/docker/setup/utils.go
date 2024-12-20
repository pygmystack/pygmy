package setup

// Unique will return a slice with duplicates
// removed. It performs a similar function to
// the linux program `uniq`
func Unique(stringSlice []string) []string {
	m := make(map[string]bool)
	for _, item := range stringSlice {
		if _, ok := m[item]; !ok {
			m[item] = true
		}
	}

	var result []string
	for item := range m {
		result = append(result, item)
	}
	return result
}
