package utils

import (
	"log"
)

func FatalCall[T any](fn func() (T, error)) T {
	val, err := fn()
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func FatalCallErrorSupplier(fn func() error) {
	err := fn()
	if err != nil {
		log.Fatal(err)
	}
}
