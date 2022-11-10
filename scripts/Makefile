## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://suva.sh/posts/well-documented-makefiles/

SHELL=bash
TARGET_CFG?=default.server.yml
TARGET_NAME?=

SOURCE_ROOT?=$(realpath $(firstword $(subst :, ,${GOPATH}))/src/github.com/yunling101/cayo)
VERSION?=$(shell cat ${SOURCE_ROOT}/.changelog | head -n 1)
BUILD_DIR?=$(SOURCE_ROOT)/cayo-${VERSION}

GOBUILD_FLAGS=-ldflags "-w -s"
GO = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go

# building
${TARGET_NAME}:
	@echo -e "\033[34mbuilding start ${TARGET_NAME}... \033[0m"
	@[ -d ${BUILD_DIR}/certs ]  || mkdir -p ${BUILD_DIR}/certs
	@cp -r ${SOURCE_ROOT}/certs/* ${BUILD_DIR}/certs
	
	@[ -d ${BUILD_DIR}/template ]  || mkdir -p ${BUILD_DIR}/template
	@cp -r ${SOURCE_ROOT}/template/* ${BUILD_DIR}/template
	
	@[ -f ${BUILD_DIR}/LICENSE ]  || cp ${SOURCE_ROOT}/LICENSE ${BUILD_DIR}
	@[ -f ${BUILD_DIR}/INSTALL ]  || cp ${SOURCE_ROOT}/INSTALL ${BUILD_DIR}

	@[ -d ${BUILD_DIR}/bin ]  || mkdir -p ${BUILD_DIR}/bin
	@[ -d ${BUILD_DIR}/config ] || mkdir -p ${BUILD_DIR}/config
	@[ -f $(BUILD_DIR)/config/${TARGET_CFG} ] || cp -f $(SOURCE_ROOT)/config/${TARGET_CFG} $(BUILD_DIR)/config/${TARGET_CFG}
	@$(GO) build $(GOBUILD_FLAGS) -o ${BUILD_DIR}/bin/${TARGET_NAME}
	@echo
	@echo ${BUILD_DIR}/bin/${TARGET_NAME}