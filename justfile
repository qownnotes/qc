# Use `just <recipe>` to run a recipe
# https://just.systems/man/en/

# By default, run the `--list` command
default:
    @just --list

# Aliases

alias fmt := format

# Download dependencies
[group('build')]
dep:
    go mod download

# Update dependencies
[group('build')]
update:
    go get -u
    go mod tidy

# Build the project
[group('build')]
build:
    go build -o qc main.go

# Execute the built binary
[group('build')]
exec +ARGS='':
    ./qc exec {{ARGS}}

# Build and execute in one command
[group('build')]
build-exec +ARGS='': build
   just exec {{ARGS}}

# Install the project
[group('build')]
install:
    go install

# Run tests
[group('test')]
test:
    go test ./...

# Run go vet
[group('test')]
vet:
    go vet

# Build using nix
[group('nix')]
nix-build:
    nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })'

# Force build using nix
[group('nix')]
nix-build-force:
    nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })' --check

# Format all justfiles
[group('linter')]
just-format:
    #!/usr/bin/env bash
    # Find all files named "justfile" recursively and run just --fmt --unstable on them
    find . -type f -name "justfile" -print0 | while IFS= read -r -d '' file; do
        echo "Formatting $file"
        just --fmt --unstable -f "$file"
    done

# Format all files
[group('linter')]
format args='':
    nix-shell -p treefmt go nodePackages.prettier shfmt nixfmt-rfc-style statix taplo --run "treefmt {{ args }}"
