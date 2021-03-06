package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	configFile    string
	line          string
	err           error
	linesChan     chan string
	lines         = []string{}
	mux           = &sync.Mutex{}
	cfg           *config
	logFile       *os.File
	writeOutAsCsv bool
	writeOutAsRaw bool
)

type config struct {
	Columns       []string `json:"columns"`
	WriteInterval int64    `json:"write_interval"`
	Commad        string   `json:"command"`
	CommadArgs    []string `json:"command_args"`
	OutDir        string   `json:"out_dir"`
	LogDir        string   `json:"log_dir"`
	CsvDelimiter  string   `json:"csv_delimiter"`
	OutFormats    []string `json:"out_formats"`
}

func init() {
	var (
		help bool
	)
	flag.StringVar(&configFile, "cfg", "config.json", "Configuration file")
	flag.BoolVar(&help, "h", false, "Print help")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	cfg = readConfig()
	logFileName := path.Join(cfg.LogDir, "mosaic_go.log")
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    20, // megabytes
		MaxBackups: 5,
	})

	for _, format := range cfg.OutFormats {
		switch strings.ToLower(format) {
		case "csv":
			writeOutAsCsv = true
		case "raw":
			writeOutAsRaw = true
		}
	}
}

func readConfig() *config {
	var (
		err     error
		cfgBuff []byte
	)
	if cfgBuff, err = ioutil.ReadFile(configFile); err != nil {
		fmt.Printf("Error opening config file. %s, %s\n", configFile, err)
		os.Exit(1)
	}
	cfg := &config{}
	if err = json.Unmarshal(cfgBuff, cfg); err != nil {
		fmt.Printf("Error reading configuration from JSON. %s, %s\n", configFile, err)
		os.Exit(1)
	}
	return cfg
}

func writeCsv(lines [][]string, fileName string, cfg *config) (err error) {
	if !writeOutAsCsv {
		return
	}
	fileName = fmt.Sprintf("%s.%s", fileName, "csv")
	var f *os.File
	if f, err = os.Create(fileName); err != nil {
		return
	}
	defer f.Close()
	csvWriter := csv.NewWriter(f)
	csvWriter.Comma = rune(cfg.CsvDelimiter[0])
	if len(lines) == 0 {
		return
	}
	for _, line := range lines {
		if err = csvWriter.Write(line); err != nil {
			return
		}
	}
	csvWriter.Flush()
	log.Printf("%d records stored in %s", len(lines), fileName)
	return
}

func writeRawOut(lines []string, fileName string, cfg *config) (err error) {
	if !writeOutAsRaw {
		return
	}
	fileName = fmt.Sprintf("%s.%s", fileName, "json")
	err = ioutil.WriteFile(fileName, []byte(strings.Join(lines, "\n")), os.ModePerm)
	log.Printf("%d records stored in %s", len(lines), fileName)
	return
}

func writeChunk(lines []string, cfg *config) (err error) {
	if len(lines) == 0 {
		return
	}
	fileName := fmt.Sprintf("mosaic_%s%d", time.Now().Format("20060102150405"), time.Now().Nanosecond())
	fileName = path.Join(cfg.OutDir, fileName)
	csvLines := [][]string{}
	for _, line := range lines {
		csvLine := []string{}
		lineMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(line), &lineMap)
		if err != nil {
			log.Printf("Error in JSON line: %v: %s\n", err, line)
			continue
		}
		for _, column := range cfg.Columns {
			colVal := lineMap[column]
			colValue := ""
			switch colVal.(type) {
			case string:
				colValue = colVal.(string)
			case nil:
				colValue = ""
			default:
				colValue = fmt.Sprintf("%v", colVal)
			}

			if column == "VIOLATION_DATE" && colValue == "none" {
				colValue = ""
			}
			if column == "FAULT_IMPACT_TYPE_ID" {
				colValue = strings.ToLower(colValue)

				switch colValue {
				case "nsa", "non-service affecting":
					colValue = "NSA"
				case "sa", "service affecting":
					colValue = "SA"
				default:
					colValue = "UNK"
				}
			}
			csvLine = append(csvLine, colValue)
		}
		csvLines = append(csvLines, csvLine)
	}

	if err = writeCsv(csvLines, fileName, cfg); err != nil {
		return
	}

	err = writeRawOut(lines, fileName, cfg)
	return
}
func main() {
	defer logFile.Close()
	cfg := readConfig()
	log.Printf("Run command %s %s\n", cfg.Commad, strings.Join(cfg.CommadArgs, " "))
	cmd := exec.Command(cfg.Commad, cfg.CommadArgs...)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	linesChan = make(chan string)
	go func() {
		for {
			select {
			case line := <-linesChan:
				mux.Lock()
				lines = append(lines, line)
				mux.Unlock()
			}
		}
	}()

	scanner := bufio.NewScanner(stdout)
	go func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			m := scanner.Text()
			linesChan <- m
		}
	}(scanner)

	ticker := time.NewTicker(time.Second * time.Duration(cfg.WriteInterval))
	go func(ticker *time.Ticker, cfg *config) {
		var err error
		for _ = range ticker.C {
			mux.Lock()
			if err = writeChunk(lines, cfg); err != nil {
				log.Printf("Error write file: %s\n", err)
			}

			lines = lines[:0]
			mux.Unlock()
		}
	}(ticker, cfg)
	if err = cmd.Wait(); err != nil {
		log.Printf("Error for commad %s. %s\n", cfg.Commad, err)
	}
}
