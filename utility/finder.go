package utility

import "sort"

func IsIn(dataList []string, value string) bool {
	sort.Strings(dataList)
	return sort.SearchStrings(dataList, value) != len(dataList)
}
