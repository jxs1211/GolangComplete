.PHONY: generate
generate:
	go generate ./...

.PHONY: build-mode
build-mode:
	go build -buildmode=plugin -o ./plugin/plugin.so ./plugin/plugin.go

.PHONY: build
build:
	go build -o main ./main.go

# .PHONY: clean
# 	rm -f ./plugin/plugin.go 

.PHONY: all
all: generate build-mode build
	echo done