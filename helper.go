package main

func GetMap(segments []string) map[string]bool {
	result := make(map[string]bool)
	for _, segment := range segments {
		if _, uniqueStr := result[segment]; !uniqueStr {
			result[segment] = true
		}
	}
	return result
}