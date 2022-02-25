package main

import (
	"github.com/Masterminds/sprig/v3"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"text/template"
)

type Args struct {
	Cmd            string
	TemplateString string `positional-arg-name:"template-string"`
	InputString    string `positional-arg-name:"input-string" default:""`
}

type Options struct {
	Args         `positional-args:"true"`
	TemplateFile string `long:"template" short:"t" description:"Template filename"`
	InputFile    string `long:"input" short:"i" description:"Input filename"`
}

func main() {
	var options Options

	_, err := flags.ParseArgs(&options, os.Args)
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// Trying source in this order: string, file, stdin, empty string.
	var input []byte
	if options.InputString == "" {
		var source *os.File
		if options.InputFile != "" {
			if source, err = os.Open(options.InputFile); err != nil {
				panic(err)
			}
		} else {
			stat, _ := os.Stdin.Stat()
			if stat.Mode()&os.ModeNamedPipe != 0 {
				source = os.Stdin
			}
		}
		if source != nil {
			defer func() {
				_ = source.Close()
			}()
			input, err = ioutil.ReadAll(source)
			if err != nil {
				panic(err)
			}
		} else {
			input = []byte{}
		}
	} else {
		input = []byte(options.InputString)
	}
	var data interface{}
	if err = yaml.Unmarshal(input, &data); err != nil {
		panic(err)
	}

	templateString := options.TemplateString
	if templateString == "" {
		if options.TemplateFile != "" {
			var source *os.File
			if source, err = os.Open(options.TemplateFile); err != nil {
				panic(err)
			}
			defer func() {
				_ = source.Close()
			}()
			inputBytes, err := ioutil.ReadAll(source)
			if err != nil {
				panic(err)
			}
			templateString = string(inputBytes)
		}
	}

	t := template.Must(
		template.New("go-template").Funcs(sprig.TxtFuncMap()).Parse(templateString))

	if err = t.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
