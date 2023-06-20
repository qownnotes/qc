.PHONY: \
	dep \
	install \
	build \
	vet \
	test

dep:
	go mod download

build: main.go
	go build -o qc $<

exec: main.go
	./qc exec

build-exec: main.go
	make build && make exec

install: main.go
	go install

test:
	go test ./...

vet:
	go vet

nix-build:
	nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })'

nix-build-force:
	nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })' --check
