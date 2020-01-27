package handler

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

type (
	random struct{}
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*()"
)

var (
	customRand *rand.Rand
)

func init() {
	customRand = rand.New(&random{})
}

func (r *random) Int63() int64 {
	var b [8]byte
	_, err := cryptoRand.Read(b[:])
	if err != nil {
		panic("error read random bytes")
	}

	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

func (r *random) Seed(_ int64) {}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[customRand.Intn(len(letterBytes))]
	}
	return string(b)
}
