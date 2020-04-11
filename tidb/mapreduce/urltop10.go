package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// URLTop10 .
func URLTop10(nWorkers int) RoundsArgs {
	// YOUR CODE HERE :)
	// And don't forget to document your idea.
	var args RoundsArgs
	// round 1: do url count
	args = append(args, RoundArgs{
		MapFunc:    URLCountMap,
		ReduceFunc: URLCountReduce,
		NReduce:    nWorkers,
	})
	// round 2: sort and get the 10 most frequent URLs
	args = append(args, RoundArgs{
		MapFunc:    URLTop10Map,
		ReduceFunc: URLTop10Reduce,
		NReduce:    1,
	})
	return args
}

// URLCountMap is the map function in the first round
func URLCountMap(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	kvs := make([]KeyValue, 0, len(lines))
	m := make(map[string]int)
	for i := range lines {
		l := strings.TrimSpace(lines[i])
		if len(l) == 0 {
			continue
		}
		m[l]++
	}

	for k, v := range m {
		kvs = append(kvs, KeyValue{Key: k, Value: strconv.Itoa(v)})
	}
	return kvs
}

// URLCountReduce is the reduce function in the first round
func URLCountReduce(key string, values []string) string {
	count := 0
	for i := range values {
		n, _ := strconv.Atoi(values[i])
		count += n
	}
	return fmt.Sprintf("%s %d\n", key, count)
}

// URLTop10Map is the map function in the second round
func URLTop10Map(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	urls := make(UrlHeap, 0, 11)
	for i := range lines {
		if len(lines[i]) == 0 {
			continue
		}
		temp := urlCount{}
		_, e := fmt.Sscanf(lines[i], "%s %d", &temp.url, &temp.cnt)
		if e != nil {
			panic(len(lines[i]))
		}

		heap.Push(&urls, temp)
		if len(urls) > 10 {
			heap.Pop(&urls)
		}
	}

	kvs := make([]KeyValue, urls.Len())
	for i := range urls {
		kvs[i].Value = urls[i].url + " " + strconv.Itoa(urls[i].cnt)
	}

	return kvs
}

// URLTop10Reduce is the reduce function in the second round
func URLTop10Reduce(key string, values []string) string {
	urls := make(UrlHeap, 0, 11)
	for _, v := range values {
		v := strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		tmp := strings.Split(v, " ")
		n, err := strconv.Atoi(tmp[1])
		if err != nil {
			panic(err)
		}
		heap.Push(&urls, urlCount{tmp[0], n})
		if urls.Len() > 10 {
			heap.Pop(&urls)
		}
	}

	sort.Slice(urls, func(i, j int) bool {
		if urls[i].cnt == urls[j].cnt {
			return urls[i].url < urls[j].url
		}
		return urls[i].cnt > urls[j].cnt
	})
	buf := new(bytes.Buffer)
	for i := range urls {
		fmt.Fprintf(buf, "%s: %d\n", urls[i].url, urls[i].cnt)
	}
	return buf.String()
}
