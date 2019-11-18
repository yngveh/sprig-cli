package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
)

var (
	datafileFlag = flag.String("data", "", "Datafile")
	tmplFlag     = flag.String("tmpl", "", "Template")

	tmpl []byte
	data map[string]interface{}
)

func init() {
	flag.Parse()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeNamedPipe == 0) && *tmplFlag == "" {
		panic("No template ")
	}

	var (
		source *os.File
	)

	if *tmplFlag == "-" {
		source = os.Stdin
		defer source.Close()
	} else {
		files := strings.Split(*tmplFlag, ",")
		sources := make([]*os.File, len(files))
		for idx, file := range files {
			source, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			sources[idx] = source
			defer source.Close()
		}
		for _, source := range sources {
			aTmpl, err := ioutil.ReadAll(source)
			if err != nil {
				panic(err)
			}
			tmpl = append(tmpl, aTmpl...)
		}
	}

	if *datafileFlag != "" {
		files := strings.Split(*datafileFlag, ",")
		out, err := parseAll(files)
		if err != nil {
			panic(err)
		}
		data = out
	}
}

func parseAll(filepaths []string) (map[string]interface{}, error) {
	maps := make([]map[string]interface{}, len(filepaths))

	for idx, filepath := range filepaths {
		d, err := parse(filepath)
		if err != nil {
			return nil, err
		}
		maps[idx] = d
	}

	output := make(map[string]interface{})
	for _, d := range maps {
		output = merge(output, d)
	}

	return output, nil
}

func parse(filepath string) (map[string]interface{}, error) {
	var d map[string]interface{}

	dataBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dataBytes, &d)
	if err != nil {
		return nil, err
	}

	d, err = CastKeysToStrings(d)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func CastKeysToStrings(s interface{}) (map[string]interface{}, error) {
	new := map[string]interface{}{}
	switch src := s.(type) {
	case map[interface{}]interface{}:
		for k, v := range src {
			var str_k string
			switch typed_k := k.(type) {
			case string:
				str_k = typed_k
			default:
				return nil, fmt.Errorf("unexpected type of key in map: expected string, got %T: value=%v, map=%v", typed_k, typed_k, src)
			}

			casted_v, err := recursivelyStringifyMapKey(v)
			if err != nil {
				return nil, err
			}

			new[str_k] = casted_v
		}
	case map[string]interface{}:
		for k, v := range src {
			casted_v, err := recursivelyStringifyMapKey(v)
			if err != nil {
				return nil, err
			}

			new[k] = casted_v
		}
	}
	return new, nil
}

func recursivelyStringifyMapKey(v interface{}) (interface{}, error) {
	var casted_v interface{}
	switch typed_v := v.(type) {
	case map[interface{}]interface{}, map[string]interface{}:
		tmp, err := CastKeysToStrings(typed_v)
		if err != nil {
			return nil, err
		}
		casted_v = tmp
	case []interface{}:
		a := []interface{}{}
		for i := range typed_v {
			res, err := recursivelyStringifyMapKey(typed_v[i])
			if err != nil {
				return nil, err
			}
			a = append(a, res)
		}
		casted_v = a
	default:
		casted_v = typed_v
	}
	return casted_v, nil
}


// merge takes two maps and merges them. on collision b overwrites a
func merge(a, b map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range a {
		output[k] = v
	}

	for k, v := range b {
		output[k] = v
	}

	return output
}

func main() {
	t := template.Must(template.New(*tmplFlag).Funcs(sprig.TxtFuncMap()).Parse(string(tmpl)))
	if err := t.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
