package main

import (
	"log"

	zc "github.com/takanoriyanagitani/go-zip-cat"
)

func main() {
	e := zc.StdinToZipFilenamesToConcatenatedToStdout()
	if nil != e {
		log.Printf("%v\n", e)
	}
}
