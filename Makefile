SHELL := /bin/bash

DIR := ${CURDIR}

build:
	docker build -t pygmy .
	@echo "Removing binaries from previous build"
	docker run --rm -v $(DIR):/data pygmy sh -c 'rm -f /data/builds/pygmy*'
	@echo "Done"
	@echo "Copying binaries to build directory"
	docker run --rm -v $(DIR):/data pygmy sh -c 'cp pygmy* /data/builds/.'
	@echo "Done"
	@echo "Enjoy using pygmy binaries in the $(DIR)/builds directory."

clean:
	docker image rm -f pygmy
	docker image prune -f --filter="label=stage=builder"

