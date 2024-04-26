package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	pattern   string
	replace   string
	patternRE *regexp.Regexp
	usage     = `Extract, transform data from stdin to stdout
expr [pattern] -replace [replace_expression]
	pattern (string) - regexp search expression
	replace () - transform expression

	`
)

func init() {
	var (
		err error
	)
	flag.StringVar(&replace, "replace", "", "replace value")
	flag.StringVar(&pattern, "pattern", "", "pattern")
	fmt.Println(os.Args)
	flag.Parse()
	if patternRE, err = regexp.Compile(pattern); err != nil {
		fmt.Printf("Wrong RE pattern: %s, %vn", pattern, err)
		os.Exit(1)
	}
	if replace != "" {
		replace = strings.Replace(replace, "\\", "$", -1)
	}
}

func main() {
	var (
		bLine []byte
		err   error
	)
	_, err = os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)
	//var output [][]byte

	for {
		bLine, _, err = reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		if patternRE.Match(bLine) {
			if replace == "" {
				fmt.Println(string(bLine))
				continue
			}
			fmt.Println(string(patternRE.ReplaceAll(bLine, []byte(replace))))
		}
		//output = append(output, bLine)
	}

	// for j := 0; j < len(output); j++ {
	// 	fmt.Printf("%d >> %s\n", j+1, output[j])
	// }

}
