package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

const lie = "🍰" // 4 bytes

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

func BenchmarkProcess(b *testing.B) {
	for n := 0; n < b.N; n++ {
		process("./manifest.json")
	}
}
