GO_PLUGIN_BUILD2=go build -buildmode=plugin -trimpath -gcflags="all=-N -l"
GO_PLUGIN_BUILD=go build -buildmode=plugin
GO_BUILD=go build
GO_PLUGIN_OBJ_PATH=./plugin.so
GO_PLUGIN_SRC_PATH=./plugin/plugin.go
MAIN_FILE=main

.PHONY: all
all: generate build-mode build run clean
	@echo done

.PHONY: generate
generate:
	go generate ./...

.PHONY: build-mode
build-mode:
	${GO_PLUGIN_BUILD} -o ${GO_PLUGIN_OBJ_PATH} ./plugin/plugin.go

.PHONY: build
build:
	${GO_BUILD} -o ${MAIN_FILE} ./main.go

.PHONY: run program
run: 
	@./${MAIN_FILE}

.PHONY: clean
clean:
	rm -f ${GO_PLUGIN_OBJ_PATH} ${GO_PLUGIN_SRC_PATH} ./${MAIN_FILE}
