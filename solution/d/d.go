package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"container/heap"
)

var in = bufio.NewReader(os.Stdin)

func main() {

	t := readInt(in)
	for i := 0; i < t; i++ {
		solve()
	}

}

func solve() {

}

func minMax(ns ...int) (int, int) {
	m, M := ns[0], ns[0]
	for i := 1; i < len(ns); i++ {
		if ns[i] < m {
			m = ns[i]
		}
		if ns[i] > M {
			M = ns[i]
		}
	}
	return m, M
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func readInt(in *bufio.Reader) int {
	l, _ := strconv.Atoi(readLine(in))
	return l
}

func readLine(in *bufio.Reader) string {
	l, _ := in.ReadString('\n')
	return strings.TrimSpace(l)
}

func readArrString(in *bufio.Reader) []string {
	return strings.Split(readLine(in), " ")
}

func readArrInt(in *bufio.Reader) []int {
	r := readArrString(in)
	arr := make([]int, len(r))
	for i := 0; i < len(arr); i++ {
		arr[i], _ = strconv.Atoi(r[i])
	}
	return arr
}

func readArrInt64(in *bufio.Reader) []int64 {
	r := readArrString(in)
	arr := make([]int64, len(r))
	for i := 0; i < len(arr); i++ {
		arr[i], _ = strconv.ParseInt(r[i], 10, 64)
	}
	return arr
}

func write(arg ...interface{})            { fmt.Print(arg...) }
func writeLine(arg ...interface{})        { fmt.Println(arg...) }
func writeF(f string, arg ...interface{}) { fmt.Printf(f, arg...) }


type PQ []int

var _ heap.Interface = new(PQ)

func (q PQ) Len() int           { return len(q) }
func (q PQ) Less(i, j int) bool { return q[i] > q[j] } // TODO: change Less as needed by the question
func (q PQ) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *PQ) Pop() interface{} {
	arr := *q
	v := arr[len(arr)-1]
	*q = arr[:len(arr)-1]
	return v
}
func (q *PQ) Push(v interface{}) { *q = append(*q, v.(int)) }
type unionFind struct {
    Parent []int
    rank   []int
}
func NF(size int) *unionFind {
    uf := unionFind{
        Parent: make([]int, size),
        rank:   make([]int, size),
    }
    for i := range uf.Parent {
        uf.Parent[i] = i
    }
    return &uf
}

func (uf *unionFind) Find(x int) int {
    if uf.Parent[x] != x {
        uf.Parent[x] = uf.Find(uf.Parent[x])
    }
    return uf.Parent[x]
}

func (uf *unionFind) Union(x, y int) {
    px, py := uf.Find(x), uf.Find(y)
    if px == py {
        return
    }

    if uf.rank[px] < uf.rank[py] {
        uf.Parent[px] = py
    } else if uf.rank[px] > uf.rank[py] {
        uf.Parent[py] = px
    } else {
        uf.Parent[px] = py
        uf.rank[py]++
    }
}