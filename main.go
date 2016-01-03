// main
package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	runtime.Gosched()
	Routes()
}
