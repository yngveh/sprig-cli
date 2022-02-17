# Sprig command line tool

Command line for [Sprig](https://github.com/Masterminds/sprig). Works with YAML and for JSON as part of YAML.

## Install

```bash
go install github.com/ksa-real/sprig-cli@latest
```

## Usage

### Data from argument

```bash
sprig-cli -t /path/to/template -i /path/to/data.yaml
sprig-cli "{{ .a }}" "a: vvv" 
sprig-cli '{{ .a.b }}' '{ a: { b: "ccc" }}'
sprig-cli '{{ .a.b }}' -i /path/to/data.yaml
```

### Data from stdin

```bash
sprig-cli '{{ uuidv4 }}'
cat test/my-data.yaml | spring-cli "{{ .root.key1 }}" 
```

## Build
```bash
go build
```

