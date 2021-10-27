package utility

// IsIn check if the value exists in the dataList
// don't use it for large data set
// time complexity O(n) because of linear scan
func IsIn(dataList []string, value string) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}
