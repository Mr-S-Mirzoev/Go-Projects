package main

import (
	"sort"
	"strings"
)

// CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/), объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
func CombineResults(in, out chan interface{}) {
	var s string
	results := make([]string, 0, 10)

	for iface := range in {
		s = iface.(string)
		results = append(results, s)
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}
