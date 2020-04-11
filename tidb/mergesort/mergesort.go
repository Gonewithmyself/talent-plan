package main

import (
	"runtime"
	"sync/atomic"
)

// MergeSort performs the merge sort algorithm.
// Please supplement this function toMerge accomplish the home work.

// MergeSortMutiGoroutine .
func MergeSortMutiGoroutine(s []int64) {
	total := len(s)
	if total < 10000 {
		MergeSort(s)
		return
	}
	n := runtime.NumCPU()
	sz := len(s) / n
	left := len(s) % n
	type sorted struct {
		start, end int
	}
	table := []*sorted{}
	temp := make([]int64, len(s))

	var count int64
	ch := make(chan *sorted, 2)
	for i := 0; i < n; i++ {
		go func(i int) {
			start := i * sz
			end := start + sz
			if i == n-1 {
				end += left
			}
			MergeSort(s[start:end])
			ch <- &sorted{start, end}
			if atomic.AddInt64(&count, 1) == int64(n) {
				close(ch)
			}
		}(i)

	}

	exit := false
	var cur *sorted
	mergeOnce := func() (merged bool) {
		n := len(table)
		if n == 0 {
			return
		}
		cur = table[n-1]
		if cur.start == 0 && cur.end == total {
			exit = true
			return
		}

		for i := 0; i < n-1; i++ {
			toMerge := table[i]
			eq := -1
			if toMerge.end == cur.start {
				eq = toMerge.end
			} else if toMerge.start == cur.end {
				eq = toMerge.start
			}
			if eq != -1 {
				left, right := min(cur.start, toMerge.start), max(cur.end, toMerge.end)
				merge(s[left:right], 0, eq-1-left, right-1-left, temp)
				table[i].start = left
				table[i].end = right
				table = table[:n-1]
				merged = true
				break
			}
		}
		return
	}

	for item := range ch {
		table = append(table, item)
		mergeOnce()
	}

	for !exit {
		for {
			if !mergeOnce() {
				break
			}
		}
	}
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

func MergeSort(s []int64) {
	n := len(s)
	temp := make([]int64, len(s))
	for sz := 1; sz < n; sz *= 2 {
		for lo := 0; lo < n-sz; lo += 2 * sz {
			end := lo + 2*sz - 1
			if n-1 < end {
				end = n - 1
			}

			if sz < 15 {
				insertSort(s[lo : end+1])
				continue
			}
			merge(s, lo, lo+sz-1, end, temp)
		}
	}
}

func insertSort(s []int64) {
	for i := 1; i < len(s); i++ {
		temp := s[i]
		j := i
		for ; j > 0; j-- {
			if s[j-1] > temp {
				s[j] = s[j-1]
			} else {
				break
			}
		}
		s[j] = temp
	}
}

func merge(s []int64, left, mid, right int, temp []int64) {
	// temp := make([]int64, right-left+1)
	i, j, k := left, mid+1, left
	for i <= mid && j <= right {
		if s[i] <= s[j] {
			temp[k] = s[i]
			i++
		} else {
			temp[k] = s[j]
			j++
		}
		k++
	}

	for i <= mid {
		temp[k] = s[i]
		i++
		k++
	}

	for j <= right {
		temp[k] = s[j]
		j++
		k++
	}
	copy(s[left:], temp[left:right+1])
}

func doMergeSort(s []int64, left, right int, temp []int64) {
	if left < right {
		mid := (left + right) >> 1
		doMergeSort(s, left, mid, temp)
		doMergeSort(s, mid+1, right, temp)
		merge(s, left, mid, right, temp)
	}
}
