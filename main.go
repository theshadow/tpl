package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"time"
)

var (
	// Version the build version of the application
	Version string
	// Author creator and copyright holder for the tool
	Author = "Xander Guzman"

	tplArg        string
	dataArg       string
	outFlag       string
	verFlag       bool
	dataStdinFlag bool
)

func main() {
	flag.Parse()
	args := flag.Args()

	if verFlag {
		fmt.Printf("Version: %s ©%d, %s\n", Version, time.Now().Year(), Author)
		os.Exit(1)
	}

	argCnt := 2
	if dataStdinFlag {
		argCnt--
	}

	if len(args) < argCnt {
		errorf("missing arguments")
		flag.Usage()
		os.Exit(1)
	}

	tplArg = args[0]
	tplFile, err := ioutil.ReadFile(tplArg)
	if err != nil {
		errorf("unable to open template file for input")
		flag.Usage()
		os.Exit(1)
	}

	var b []byte
	if dataStdinFlag {
		b, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			errorf("unable to read from stdin: %s", err)
			os.Exit(1)
		}
	} else {
		dataArg = args[1]
		b, err = ioutil.ReadFile(dataArg)
		if err != nil {
			errorf("unable to open data file for input: %s", err)
			flag.Usage()
			os.Exit(1)
		}
	}

	tpl := template.New("input")
	tpl, err = tpl.Parse(string(tplFile))
	if err != nil {
		errorf("unable to parse template: %s", err)
		os.Exit(1)
	}

	data := map[string]interface{}{}
	if err = json.Unmarshal(b, &data); err != nil {
		errorf("unable to unmarshal data: %s", err)
		os.Exit(1)
	}

	outStream := os.Stdout
	if len(outFlag) > 0 {
		outStream, err = os.OpenFile(outFlag, os.O_RDONLY, 0644)
		if err != nil {
			errorf("unable to open output file: %s", err)
			os.Exit(1)
		}
	}

	err = tpl.Execute(outStream, data)
	if err != nil {
		errorf("unable to parse execute template: %s", err)
		os.Exit(1)
	}
}

func usage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s [options] template-file data-file\n\nParameters:\n", os.Args[0])

	flag.PrintDefaults()
}

func errorf(format string, a ...interface{}) {
	fmt.Printf(format + "\n", a...)
}

func init() {
	flag.Usage = usage
	flag.StringVar(&outFlag, "out-file", "", "output is instead written to this file instead of standard out")
	flag.BoolVar(&verFlag, "version", false, "the version of the application is output and then exists")
	flag.BoolVar(&dataStdinFlag, "data-stdin", false, "makes it so that the data is read from STDIN instead of a file specified as a command line argument")
}
