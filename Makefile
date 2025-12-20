VENV_DIR = .venv
PYTHON = $(VENV_DIR)/bin/python
PIP = $(VENV_DIR)/bin/pip
.PHONY: demo test bench

demo:
	@cd demo && go run main.go

test:
	@go test -v

bench:
	@go test -bench=. -benchmem

gen:
	@cd demo/gen_data && go run gen_path.go && go run gen_data.go
