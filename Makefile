BUILD_DIR?=$(shell pwd)

# building
.PHONY:default
default:
	@cd ${BUILD_DIR}/cmd/cayoserver/ && make
