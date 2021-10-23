package utility

import "sort"

// IsIn find value in dataList.
// dataList must be sorted
func IsIn(dataList []string, value string) bool {
	return sort.SearchStrings(dataList, value) != len(dataList)
}
