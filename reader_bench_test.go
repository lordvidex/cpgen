package cpgen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func readArrIntFmt(in io.Reader) []int {
	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	return arr
}

// current impl
func readLine(in *bufio.Reader) string {
	line, _ := in.ReadString('\n')
	return strings.TrimSpace(line)
}
func readArrString(in *bufio.Reader) []string {
	return strings.Split(readLine(in), " ")
}
func readArrInt(in *bufio.Reader) []int {
	r := readArrString(in)
	arr := make([]int, len(r))
	for i := 0; i < len(r); i++ {
		arr[i], _ = strconv.Atoi(r[i])
	}
	return arr
}

func BenchmarkReadArrInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		file, err := os.Open("test.input")
		if err != nil {
			b.Fatal(err)
		}
		in := bufio.NewReader(file)
		b.StartTimer()
		readArrInt(in)
	}
}

func BenchmarkReadArrIntFmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		file, err := os.Open("test.input")
		if err != nil {
			b.Fatal(err)
		}
		in := bufio.NewReader(file)
		b.StartTimer()
		readArrIntFmt(in)
	}
}
