demo:
	@cd demo && go run main.go
test:
	@go test -v
bench:
	@go test -bench=. -benchmem