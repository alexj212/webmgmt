export PROJ_PATH=github.com/alexj212/webmgmt


export DATE := $(shell date +%Y.%m.%d-%H%M)
export LATEST_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
export BRANCH := $(shell git branch |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export BUILT_ON_IP := $(shell [ $$(uname) = Linux ] && hostname -i || hostname )
export BIN_DIR=./bin
export PACKR2_EXECUTABLE := $(shell command -v packr2  2> /dev/null)

export BUILT_ON_OS=$(shell uname -a)
ifeq ($(BRANCH),)
BRANCH := master
endif

export COMMIT_CNT := $(shell git rev-list HEAD | wc -l | sed 's/ //g' )
export BUILD_NUMBER := ${BRANCH}-${COMMIT_CNT}
export COMPILE_LDFLAGS=-s -X "main.DATE=${DATE}" \
                          -X "main.LATEST_COMMIT=${LATEST_COMMIT}" \
                          -X "main.BUILD_NUMBER=${BUILD_NUMBER}" \
                          -X "main.BUILT_ON_IP=${BUILT_ON_IP}" \
                          -X "main.BUILT_ON_OS=${BUILT_ON_OS}"



build_info: check_prereq ## Build the container
	@echo ''
	@echo '---------------------------------------------------------'
	@echo 'BUILT_ON_IP       $(BUILT_ON_IP)'
	@echo 'BUILT_ON_OS       $(BUILT_ON_OS)'
	@echo 'DATE              $(DATE)'
	@echo 'LATEST_COMMIT     $(LATEST_COMMIT)'
	@echo 'BRANCH            $(BRANCH)'
	@echo 'COMMIT_CNT        $(COMMIT_CNT)'
	@echo 'BUILD_NUMBER      $(BUILD_NUMBER)'
	@echo 'COMPILE_LDFLAGS   $(COMPILE_LDFLAGS)'
	@echo 'PATH              $(PATH)'
	@echo 'PACKR2_EXECUTABLE $(PACKR2_EXECUTABLE)'
	@echo '---------------------------------------------------------'
	@echo ''


####################################################################################################################
##
## help for each task - https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
##
####################################################################################################################
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help



####################################################################################################################
##
## Build of binaries
##
####################################################################################################################
all: example test ## build example and run tests

binaries: example ## build binaries in bin dir

create_dir:
	@mkdir -p $(BIN_DIR)
	@rm -f $(BIN_DIR)/web
	@ln -s ../web $(BIN_DIR)/web

check_prereq: create_dir
ifndef PACKR2_EXECUTABLE
	go get -u github.com/gobuffalo/packr/v2/packr2
endif
	$(warning "found packr2")



build_app: create_dir
		packr2 build -o $(BIN_DIR)/$(BIN_NAME) -a -ldflags '$(COMPILE_LDFLAGS)' $(APP_PATH)


example: build_info ## build example binary in bin dir
	@echo "build 1"
	make BIN_NAME=example APP_PATH=github.com/alexj212/webmgmt/example build_app
	@echo ''
	@echo ''



####################################################################################################################
##
## Cleanup of binaries
##
####################################################################################################################

clean_binaries: clean_example  ## clean all binaries in bin dir


clean_binary: ## clean binary in bin dir
	rm -f $(BIN_DIR)/$(BIN_NAME)

clean_example: ## clean example
	make BIN_NAME=example clean_binary



test: ## run tests
	go test $(PROJ_PATH)/_test/

fmt: ## run fmt
	go fmt $(PROJ_PATH)/...



