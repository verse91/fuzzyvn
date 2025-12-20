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
	python3 -m venv $(VENV_DIR)
	$(PIP) install datasets pandas
	$(PYTHON) demo/gen_data/down_data.py
	go run demo/gen_data/gen_data.go
clean:
	rm -rf $(VENV_DIR)
	rm -f demo/gen_data/test_paths_100k.txt
