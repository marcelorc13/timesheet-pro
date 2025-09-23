.PHONY: run
run:
	@go tool templ generate
	@go run cmd/main.go


.PHONY: setup
setup:
	go install -ldflags="-s -w" -tags="no_libsql no_mssql no_vertica no_clickhouse no_mysql no_sqlite3 no_ydb" github.com/pressly/goose/v3/cmd/goose@latest
	go get -tool github.com/a-h/templ/cmd/templ@latest

.PHONY: tidy
tidy:
	go mod tidy	
	go mod verify
	go fmt ./...

.PHONY: migrations/new
migrations/new: 
	goose create ${name} sql


.PHONY: migrations/up
migrations/up: 
	goose up 
