.PHONY: build build-all clean

build-all:
	export GOOS=linux
	export GOARCH=amd64
	ls cmd | xargs -P 1 -I{} bash -c "CC=/usr/bin/musl-gcc GOPRIVATE=\"github.com/opengovern\" GOOS=linux GOARCH=amd64 go build -tags musl -v -ldflags \"-linkmode external -extldflags '-static' -s -w\" -tags musl -o ./build/ ./cmd/{}"

build:
	./scripts/list_services > ./service-list
	cat ./service-list
	cat ./service-list | grep -v "steampipe" | xargs -P 1 -I{} bash -c "CC=/usr/bin/musl-gcc GOPRIVATE=\"github.com/opengovern\" GOOS=linux GOARCH=amd64 go build -v -ldflags \"-linkmode external -extldflags '-static' -s -w\" -tags musl -o ./build/ ./cmd/{}"

clean:
	rm -r ./build