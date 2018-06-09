package main

import (
	"os/exec"
	"flag"
	"log"
	"fmt"
	"os"
)

var (
	ppath = "."
)

func init() {
	flag.StringVar(&ppath, "path", ".", "where to do gofmt")
	flag.Parse()
}

func main() {
	gofmt := exec.Command("gofmt", "-w")
	gofmt.Args = append(gofmt.Args, ppath)
	gofmt.Stderr = os.Stderr

	out, err := gofmt.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok:", string(out))
}