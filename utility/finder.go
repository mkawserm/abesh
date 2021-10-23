package utility

import "sort"

// IsIn finds value in the dataList.
// dataList must be sorted
func IsIn(dataList []string, value string) bool {
	return sort.SearchStrings(dataList, value) != len(dataList)
}
