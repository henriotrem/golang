package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type FooReader struct{}

func (f *FooReader) Read(b []byte) (n int, err error) {
	fmt.Println("in <- ")
	return os.Stdin.Read(b)
}

type FooWriter struct{}

func (f *FooWriter) Write(b []byte) (n int, err error) {
	fmt.Println("out -> ")
	return os.Stdout.Write(b)
}

func main() {
	ReaderWriter1()
	ReaderWriter2()
}

func ReaderWriter1() {
	var (
		reader FooReader
		writer FooWriter
	)
	b := make([]byte, 4096)

	n, err := reader.Read(b)
	if err != nil {
		log.Fatalln("Fatal error")
	}
	fmt.Printf("%d Bytes read", n)
	n, err = writer.Write(b)
	if err != nil {
		log.Fatalln("Fatal error")
	}
	fmt.Printf("%d Bytes written", n)
}

func ReaderWriter2() {
	var (
		reader FooReader
		writer FooWriter
	)
	if _, err := io.Copy(&writer, &reader); err != nil {
		log.Fatalln("Unable to read/write data")
	}
}
