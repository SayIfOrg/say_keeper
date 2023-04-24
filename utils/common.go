package utils

import (
	"sort"
)

//RemoveByIndexes removes elements from a list based on specified indexes
func RemoveByIndexes[T interface{}](list []T, indexes []int) []T {
	sort.Sort(sort.Reverse(sort.IntSlice(indexes)))

	// Remove elements from the list based on the indexes
	for _, index := range indexes {
		list = append(list[:index], list[index+1:]...)
	}

	return list
}
