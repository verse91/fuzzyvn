demo:
	@cd demo && go run main.go
test:
	@go test -v fuzzyvn_test.go
