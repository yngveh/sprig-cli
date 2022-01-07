# Sprig command line tool

Command line for [Sprig](https://github.com/Masterminds/sprig)

## Install

```bash
go install github.com/yngveh/sprig-cli@latest
```

## Usage

### Template from file

```bash
sprig-cli -tmpl /path/to/template.tpl -data /path/to/data.yaml
```

### Template from stdin

```bash
echo "{{ uuidv4 }}" | sprig-cli
echo "{{ .root.key1 }}" | spring-cli -data test/my-data.yaml
```

## Build

```bash
make
```

