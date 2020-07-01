package internal

func hasRole(target string, list []string) bool {
	for _, b := range list {
		if b == target {
			return true
		}
	}
	return false
}

func removeFromSlice(removeString string, list []string) []string {
	indexItem := 0

	for index, item := range list {
		if item == removeString {
			indexItem = index
		}
	}

	return append(list[:indexItem], list[indexItem+1:]...)
}
