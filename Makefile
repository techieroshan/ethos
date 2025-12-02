.PHONY: proto generate test build

# Generate Go code from Protocol Buffers
proto:
	@echo "Generating gRPC code from .proto files..."
	@mkdir -p api/proto/feedback api/proto/dashboard api/proto/notifications api/proto/people
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/feedback/feedback.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/dashboard/dashboard.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/notifications/notifications.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/people/people.proto
	@echo "gRPC code generation complete"

# Generate all code
generate: proto

# Run tests
test:
	go test ./... -v

# Build the application
build:
	go build ./cmd/api/main.go

