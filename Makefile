SHELL := /bin/bash

DIR := ${CURDIR}

build:
	docker build -t pygmy-go .
	@echo "Removing binaries from previous build"
	docker run --rm -v $(DIR):/data pygmy-go sh -c 'rm -f /data/builds/pygmy-go*'
	@echo "Done"
	@echo "Copying binaries to build directory"
	docker run --rm -v $(DIR):/data pygmy-go sh -c 'cp pygmy-g* /data/builds/.'
	@echo "Done"
	@echo "Enjoy using pygmy-go binaries in the $(DIR)/build directory."

clean:
	docker image rm -f pygmy-go
	docker image prune -f --filter="label=stage=builder"

