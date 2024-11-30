package main

// arrayColumn extracts a column of values from a slice of maps.
// - array: The slice of maps (like rows in a table).
// - key: The key to extract values for.
// Returns a slice of values corresponding to the specified key.
func arrayColumn[K comparable, V any](array []map[K]V, key K) []V {
	result := make([]V, 0, len(array))
	for _, row := range array {
		if value, ok := row[key]; ok {
			result = append(result, value)
		}
	}
	return result
}
