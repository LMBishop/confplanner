BINARY_NAME=confplanner

all: build

.PHONY: build
build: web
	go build -o ${BINARY_NAME} main.go

.PHONY: web
web:
	(cd web && go generate)

.PHONY: sqlc
sqlc:
	sqlc compile
