package challenge

import (
	"log"
	"path/filepath"
	"runtime"
)

// Gets the openapi specification path.
func GetSpecPath() string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	relativePath := currentDir + "/../openapi.json"

	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		log.Fatal(err)
	}

	return absolutePath
}
