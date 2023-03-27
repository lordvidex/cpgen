GOBIN=$(shell go env GOPATH)/bin

run-cli:
	go run cmd/cli/cli.go
run-tui:
	go run cmd/tui/tui.go
install-cli:
	go build -o ./bin/cpgen-cli ./cmd/cli/cli.go && \
	mv ./bin/cpgen-cli $(GOBIN)/cpgen-cli
install-tui:
	go build -o ./bin/cpgen-tui ./cmd/tui/tui.go && \
	mv ./bin/cpgen-tui $(GOBIN)/cpgen
