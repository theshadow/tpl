package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/theshadow/tpl/plugins"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
)

const (
	PluginPathEnvName = "TPL_PLUGIN_PATH"
)

var (
	// Version the build version of the application
	Version string
	// Author creator and copyright holder for the tool
	Author = "Xander Guzman"

	dataArg string
	outFlag string
	tplArg  string

	verFlag         bool
	dataStdinFlag   bool
	pluginsPathFlag string
	pluginsFlag     string
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

	// override where the plugins are searched for using an OS environment variable
	if len(os.Getenv(PluginPathEnvName)) > 0 {
		pluginsPathFlag = os.Getenv(PluginPathEnvName)
	}

	// start processing STDIN
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

	// create the template
	tpl := template.New("input")
	mgr := plugins.New()

	// register the functions for each of the plugins
	if len(pluginsFlag) > 0 {
		// turn the plugins flag into a list of paths to look for plugin files
		paths := strings.Split(pluginsFlag, ",")
		for i, p := range paths {
			paths[i] = fmt.Sprintf("%s/%s.so", pluginsPathFlag, p)
		}

		err = mgr.Load(paths...)
		if err != nil {
			errorf("unable to load plugin %s", err)
			os.Exit(1)
		}

		funcMap := make(template.FuncMap)
		for pn, p := range mgr.Plugins {
			for n, fn := range p.Functions() {
				funcMap[fmt.Sprintf("%s_%s", pn, n)] = fn.(interface{})
			}
		}

		tpl.Funcs(funcMap)
	}



	tpl, err = tpl.Parse(string(tplFile))
	if err != nil {
		errorf("unable to parse template: %s", err)
		os.Exit(1)
	}

	// marshall the template data
	data := map[string]interface{}{}
	if err = json.Unmarshal(b, &data); err != nil {
		errorf("unable to unmarshal data: %s", err)
		os.Exit(1)
	}

	// render the parsed and processed template
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
	fmt.Printf(format+"\n", a...)
}

func init() {
	flag.Usage = usage
	flag.StringVar(&outFlag, "out-file", "", "output is instead written to this file instead of standard out")
	flag.BoolVar(&verFlag, "version", false, "the version of the application is output and then exists")
	flag.BoolVar(&dataStdinFlag, "data-stdin", false, "makes it so that the data is read from STDIN instead of a file specified as a command line argument")
	flag.StringVar(&pluginsPathFlag, "plugins-path", ".", "define which path to load plugins from")
	flag.StringVar(&pluginsFlag, "plugins", "", "comma delimited name of plugins without extensions to load")
}
