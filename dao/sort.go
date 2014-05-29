package dao

import (
	"sort"
)

func sortSeries(series BySize) []Series {
	for _, serie := range series {
		sort.Sort(ByDate(serie.Data))
	}
	sort.Sort(series)
	return series
}

type BySize []Series

func (s BySize) Len() int      { return len(s) }
func (s BySize) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s BySize) Less(i, j int) bool {
	iLen := len(s[i].Data) - 1
	jLen := len(s[j].Data) - 1
	iSize := s[i].Data[iLen].Value
	jSize := s[j].Data[jLen].Value
	return iSize > jSize
}

type ByDate []Point

func (s ByDate) Len() int      { return len(s) }
func (s ByDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByDate) Less(i, j int) bool {
	return s[i].Time < s[j].Time
}
