package utility

import "sort"

func IsIn(dataList []string, value string) bool {
	return sort.SearchStrings(dataList, value) != len(dataList)
}
