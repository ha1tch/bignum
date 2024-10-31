# Makefile for benchmarking bignum

# Set the benchmark directory
BENCH_DIR := benchmark

# Set the bignum package path (adjust if needed)
BIGNUM_PKG = github.com/ha1tch/bignum

# Build the bignum package
build:
	go build $(BIGNUM_PKG)

# Clean build artifacts (without deleting benchmark file)
clean:
	go clean -testcache

# Run benchmarks and build bignum
all:
	cd $(BIGNUM_PKG) && go test -bench=. -run=^$ $(BENCH_DIR)
	go build $(BIGNUM_PKG)
