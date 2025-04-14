EXAMPLE_PORT ?= 8080
EXAMPLE_HOST ?= localhost
EXAMPLE_URL_HELLO  ?= http://$(EXAMPLE_HOST):$(EXAMPLE_PORT)/api/v1/examples/hello
EXAMPLE_URL_RANDOM  ?= http://$(EXAMPLE_HOST):$(EXAMPLE_PORT)/api/v1/examples/random

.PHONY: run-example kill-example-server build-example pretty-log

run-example:
	@echo "🧹Cleaning up any old example-server processes..."
	@$(MAKE) kill-example-server

	@echo "👉Building and running example server..."
	@$(MAKE) build-example

	@echo "👉Starting server in background..."
	@PORT=$(EXAMPLE_PORT) ./bin/example-server >server.log 2>&1 & \
	SERVER_PID=$$!; \
	echo "👉Server PID: $$SERVER_PID"; \
	sleep 2; \
	echo "👉Curling endpoint..."; \
	curl -i $(EXAMPLE_URL_HELLO); \
	curl -i $(EXAMPLE_URL_RANDOM); \
	kill $$SERVER_PID; \
	echo "👉Server stopped."

	@$(MAKE) pretty-log

kill-example-server:
	@EXISTING_PID=$$(pgrep -f "./bin/example-server") || true; \
	if [ -n "$$EXISTING_PID" ]; then \
		echo "👉Found running example-server (PID: $$EXISTING_PID). Killing it..."; \
		kill -9 $$EXISTING_PID; \
	else \
		echo "✅ No previous server running."; \
	fi

build-example:
	go build -o bin/example-server ./cmd/examples/main.go

pretty-log:
	@echo "👉Pretty-printed server logs:"
	@cat server.log | while IFS= read -r line; do echo $$line | jq .; done
