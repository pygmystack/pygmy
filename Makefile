SHELL := /bin/bash

DIR := ${CURDIR}

build:
	docker build -t pygmy-go .
	@echo "Removing binaries from previous build"
	docker run -v $(DIR):/data pygmy-go rm -f /data/builds/pygmy-go*
	@echo "Done"
	@echo "Copying binaries to build directory"
	docker run -v $(DIR):/data pygmy-go cp pygmy-go-linux-x86 /data/builds/.
	docker run -v $(DIR):/data pygmy-go cp pygmy-go-darwin /data/builds/.
	docker run -v $(DIR):/data pygmy-go cp pygmy-go.exe /data/builds/.
	@echo "Done"
	@echo "Enjoy using pygmy-go binaries in $(DIR)/data/build directory."

clean:
	docker image rm -f pygmy-go
	docker image prune -f --filter label=stage=builder

