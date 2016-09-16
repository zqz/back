# this might only work with fish shell.

.PHONY: build clean doc run testrun test get-deps

DEV_ENV   = DATABASE_URL=postgres://localhost/zqz-dev?sslmode=disable
TEST_ENV  = DATABASE_URL=postgres://localhost/zqz-test?sslmode=disable

default: build

get-deps:
		go get -t ./...

clean:
		rm -f zqz

build:
		go build -o zqz *.go

doc:
		godoc -http=:6060 -index

sudorun:
		env ${DEV_ENV} ./zqz

run:
		env $(DEV_ENV) go run -v *.go -livereload -cdn="/assets"

testrun:
		env $(TEST_ENV) ginkgo watch -r | colortest

test:
		env $(TEST_ENV) go test -v ./...
