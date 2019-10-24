DIR := $(PWD)

build:
	docker build -t pygmy-go .
	docker run -v $(DIR):/data pygmy-go cp pygmy-go-linux-x86 /data/builds/.
	docker run -v $(DIR):/data pygmy-go cp pygmy-go-darwin /data/builds/.

clean:
	docker image rm pygmy-go


