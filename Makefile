.PHONY: gen-proto
gen-proto:
	@for dir in $(shell find proto -mindepth 1 -maxdepth 1 -type d); do \
		protoc --go_out=./ --go-grpc_out=./ $$dir/*.proto; \
	done

.PHONY: test
test:
	echo "Hello, World!"
