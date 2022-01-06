package main

import (
	"flag"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

func main() {

	datafileFlag := flag.String("data", "", "Datafile")
	tmplFlag := flag.String("tmpl", "", "Template")

	flag.Parse()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeNamedPipe == 0) && *tmplFlag == "" {
		panic("No template ")
	}

	var (
		source *os.File
		err    error
	)

	if *tmplFlag == "" {
		source = os.Stdin
	} else {
		source, err = os.Open(*tmplFlag)
		if err != nil {
			panic(err)
		}
	}

	defer func() {
		_ = source.Close()
	}()

	tmpl, err := ioutil.ReadAll(source)
	if err != nil {
		panic(err)
	}

	var data map[interface{}]interface{}

	if *datafileFlag != "" {

		dataBytes, err := ioutil.ReadFile(*datafileFlag)
		if err != nil {
			panic(err)
		}

		if err = yaml.Unmarshal(dataBytes, &data); err != nil {
			panic(err)
		}
	}

	t := template.Must(
		template.New(*tmplFlag).Funcs(sprig.TxtFuncMap()).Parse(string(tmpl)))

	if err = t.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
