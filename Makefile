SHELL := /bin/bash

ENGINE_DIR := engine
SHARED_DIR := $(ENGINE_DIR)/cmd/shared
BUILD_DIR  := build
MAIN_FILE  := cmd/axiom/main.go
CONFIG     := cmd/axiom/initial_config.ax

SO         := $(BUILD_DIR)/axiom.so
HEADER     := $(BUILD_DIR)/axiom.h
TEST_BIN   := $(BUILD_DIR)/test
AXIOM_BIN  := $(BUILD_DIR)/axiom

.PHONY: all run build-shared test clean

all: build-shared test

run:
	cd $(ENGINE_DIR) && go run $(MAIN_FILE) $(CONFIG)

build-shared:
	mkdir -p $(BUILD_DIR)
	cd $(ENGINE_DIR) && go build -buildmode=c-shared -o ../$(SO) ./cmd/shared
	cp $(ENGINE_DIR)/cmd/shared/*.h $(HEADER) 2>/dev/null || true

test: build-shared
	gcc test.c $(SO) -o $(TEST_BIN) -Wl,-rpath,$(abspath $(BUILD_DIR))
	DYLD_LIBRARY_PATH=$(abspath $(BUILD_DIR)) $(TEST_BIN)

clean:
	rm -rf $(BUILD_DIR)
