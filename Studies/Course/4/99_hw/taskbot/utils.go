package main

func sliceIndex(slice []string, pattern string) int {
	for i, item := range slice {
		if item == pattern {
			return i
		}
	}
	return -1
}
