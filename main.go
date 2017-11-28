package main

import (
	"flag"
	"os"
	"io/ioutil"
	"text/template"
	"gopkg.in/yaml.v2"
	"github.com/Masterminds/sprig"
)

var (
	datafileFlag = flag.String("data", "", "Datafile")
	tmplFlag     = flag.String("tmpl", "", "Template")

	tmpl []byte
	data map[interface{}]interface{}
)

func init() {

	flag.Parse()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeNamedPipe == 0) && *tmplFlag == "" {
		panic("No template ")
	}

	var (
		source *os.File
		err    error
	)

	if *tmplFlag == "-" {
		source = os.Stdin
	} else {
		source, err = os.Open(*tmplFlag)
		if err != nil {
			panic(err)
		}
	}
	defer source.Close()

	tmpl, err = ioutil.ReadAll(source)

	if *datafileFlag != "" {

		dataBytes, err := ioutil.ReadFile(*datafileFlag)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(dataBytes, &data)
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	t := template.Must(template.New(*tmplFlag).Funcs(sprig.TxtFuncMap()).Parse(string(tmpl)))
	if err := t.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
