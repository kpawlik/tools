package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var (
	sleep int64
)

func init() {
	flag.Int64Var(&sleep, "sleep", 0, "")
	flag.Parse()
}
func main() {

	fileName := flag.Args()[0]
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	dir, name := path.Split(fileName)
	baseName := strings.Split(name, ".")[0]
	name = fmt.Sprintf("%s.csv", baseName)
	fileName = path.Join(dir, name)
	time.Sleep(time.Second * time.Duration(sleep))
	ioutil.WriteFile(fileName, fileContent, os.ModePerm)
}
