package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const lie = "🍰" // 4 bytes

type writer struct {
	Dst  io.Writer
	Rate float64
}

func (w writer) Write(buf []byte) (n int, err error) {
	if len(buf) == 0 {
		return 0, nil
	}

	time.Sleep(time.Duration(float64(len(buf)) / w.Rate * float64(time.Second)))

	return w.Dst.Write(buf)
}

func prepareFile(i int) {
	out := strings.Repeat(lie, i*1024)
	ioutil.WriteFile("./testfile_"+strconv.Itoa(i), []byte(out), 0644)
}

func deleteFile(i int) {
	os.Remove("./testfile_" + strconv.Itoa(i))
}

func benchmarkProcess(i int, b *testing.B) {
	f, _ := os.Create("./testfile_" + strconv.Itoa(i))
	z, _ := os.Open("/dev/zero")
	defer z.Close()
	defer f.Close()
	w := writer{
		Dst:  f,
		Rate: 5 * 1024 * 1024, // ~5mib/s
	}
	go io.CopyN(w, z, int64(i)*1024)

	for n := 0; n < b.N; n++ {
		process("./testfile_" + strconv.Itoa(i))
	}
	deleteFile(i)
}

func BenchmarkProcess4kbytes(b *testing.B)    { benchmarkProcess(4, b) }
func BenchmarkProcess40kbytes(b *testing.B)   { benchmarkProcess(40, b) }
func BenchmarkProcess120kbytes(b *testing.B)  { benchmarkProcess(120, b) }
func BenchmarkProcess400kbytes(b *testing.B)  { benchmarkProcess(400, b) }
func BenchmarkProcess4000kbytes(b *testing.B) { benchmarkProcess(4000, b) }
