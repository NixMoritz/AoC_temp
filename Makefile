ENTRYPOINT=./cmd/aoc/aoc.go
BINARY_NAME=aoc

GEN_ENTRYPOINT=./cmd/generator/generator.go
GEN_BINARY_NAME=generator

BENCH=

YEAR=
DAY=

run: build
	./dist/$(GEN_BINARY_NAME) $(YEAR) $(DAY)
	./dist/$(BINARY_NAME) $(YEAR) $(DAY)

dev:
	go run $(ENTRYPOINT_NAME) $(YEAR) $(DAY)


all: tests build

build: build-scraper build-generator

build-scraper:
	go build -ldflags='-s -w' -o dist/$(BINARY_NAME) $(ENTRYPOINT)

build-generator:
	go build -ldflags='-s -w' -o dist/$(GEN_BINARY_NAME) $(GEN_ENTRYPOINT)

tests:
	go test ./...

bench:
	go test ./... -bench=. -benchtime 3s -run=^\# -cpu=1,20

bench-prof:
	go test . -bench=${BENCH} -benchtime 3s -run=^\# -cpu=20 -cpuprofile ./tmp/$(subst /,-,$(BENCH))_cpu.prof -memprofile ./tmp/$(subst /,-,$(BENCH))_mem.prof -o ./tmp/$(subst /,-,$(BENCH)).test


clean:
	rm -f dist/*
	go mod tidy
	go clean

vet:
	go vet

lint:
	golangci-lint run --enable-all
