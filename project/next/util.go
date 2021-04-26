package main

import "math/rand"

const base = "ABCDEFGHIJKLMNOP"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = base[rand.Intn(len(base))]
	}
	return string(b)
}
