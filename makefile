.PHONY: build doc run test get-deps

DEV_USER  = 'dylan'
DEV_DB    = 'zqz-dev'
ifndef TEST_DB
	TEST_DB   = 'zqz-test'
endif
ifndef TEST_USER
	TEST_USER = 'dylan'
endif


DEV_ENV   = DATABASE_NAME=$(DEV_DB) DATABASE_USER=$(DEV_USER)
TEST_ENV  = DATABASE_NAME=$(TEST_DB) DATABASE_USER=$(TEST_USER)

default: build

get-deps:
	  go get -t ./...

build:
		go build -v ./... -o zqz

doc:
	  godoc -http=:6060 -index

run:
	  env $(DEV_ENV) go run -v *.go

test:
	  env $(TEST_ENV) go test -v ./...
	  # env $(TEST_ENV) (go-test)
