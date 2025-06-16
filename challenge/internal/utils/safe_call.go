package utils

import (
	"log"
)

func SafeCall[T any](fn func() (T, error)) T {
	val, err := fn()
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func SafeCallErrorSupplier(fn func() error) {
	err := fn()
	if err != nil {
		log.Fatal(err)
	}
}
