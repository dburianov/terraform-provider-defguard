.PHONY: build test clean

build:
	go build -o terraform-provider-defguard

test:
	go test ./...

clean:
	rm -f terraform-provider-defguard