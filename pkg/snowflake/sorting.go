package snowflake

import (
	"sort"
)

func sortStrings(strs map[string]string) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStringList(strs map[string][]string) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStringsInt(strs map[string]int) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStringsFloat(strs map[string]float64) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}

func sortStringsBool(strs map[string]bool) []string {
	sortedStringProperties := []string{}
	for k := range strs {
		sortedStringProperties = append(sortedStringProperties, k)
	}
	sort.Strings(sortedStringProperties)
	return sortedStringProperties
}
