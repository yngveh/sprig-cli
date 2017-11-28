
APP := sprig-cli

.PHONY: build
build:
	@go build -o bin/$(APP) main.go

.PHONY: run
run: build
#	@echo "Testing" | bin/sprig-cli -tmpl -
#	@bin/sprig-cli -tmpl test/my.tpl
	@bin/sprig-cli -tmpl test/my.tpl -data test/my-data.yaml

