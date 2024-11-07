package lib

func MyInclude(slice []string, value string) bool {
	for _, str := range slice {
		if str == value {
			return true
		}
	}
	return false
}
