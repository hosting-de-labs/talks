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
	out := strings.Repeat(lie, i)
	ioutil.WriteFile("./testfile_"+strconv.Itoa(i), []byte(out), 0644)
}

func deleteFile(i int) {
	os.Remove("./testfile_" + strconv.Itoa(i))
}

func benchmarkChecksum(i int, b *testing.B) {
	prepareFile(i)
	for n := 0; n < b.N; n++ {
		checksum("./testfile_" + strconv.Itoa(i))
	}
	deleteFile(i)
}

func BenchmarkChecksum4bytes(b *testing.B)       { benchmarkChecksum(1, b) }
func BenchmarkChecksum40bytes(b *testing.B)      { benchmarkChecksum(10, b) }
func BenchmarkChecksum120bytes(b *testing.B)     { benchmarkChecksum(30, b) }
func BenchmarkChecksum400bytes(b *testing.B)     { benchmarkChecksum(100, b) }
func BenchmarkChecksum4000bytes(b *testing.B)    { benchmarkChecksum(1000, b) }
func BenchmarkChecksum400000bytes(b *testing.B)  { benchmarkChecksum(100000, b) }
func BenchmarkChecksum4000000bytes(b *testing.B) { benchmarkChecksum(1000000, b) }

func benchmarkProcess(i int, b *testing.B) {
	f, _ := os.Create("./testfile_" + strconv.Itoa(i))
	z, _ := os.Open("/dev/zero")
	defer z.Close()
	defer f.Close()
	w := writer{
		Dst:  f,
		Rate: 5 * 1024, // 5kib/s
	}
	go io.CopyN(w, z, int64(i))

	for n := 0; n < b.N; n++ {
		process("./testfile_" + strconv.Itoa(i))
	}
	deleteFile(i)
}

func BenchmarkProcess4bytes(b *testing.B)       { benchmarkProcess(4, b) }
func BenchmarkProcess40bytes(b *testing.B)      { benchmarkProcess(40, b) }
func BenchmarkProcess120bytes(b *testing.B)     { benchmarkProcess(120, b) }
func BenchmarkProcess400bytes(b *testing.B)     { benchmarkProcess(400, b) }
func BenchmarkProcess4000bytes(b *testing.B)    { benchmarkProcess(4000, b) }
func BenchmarkProcess400000bytes(b *testing.B)  { benchmarkProcess(400000, b) }
func BenchmarkProcess4000000bytes(b *testing.B) { benchmarkProcess(4000000, b) }
