package main

import (
	"os/exec"
	"flag"
	"log"
	"fmt"
	"os"
	"io"
)

var (
	ppath = "."
)

func init() {
	flag.StringVar(&ppath, "path", ".", "where to do gofmt")
	flag.Parse()
}

func runGofmt(ppath string, errW io.Writer) (out []byte, err error) {
	gofmt := exec.Command("gofmt", "-w", ppath)
	gofmt.Stderr = errW
	return gofmt.Output()
}

func main() {
	out, err := runGofmt(ppath, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok:", string(out))
}