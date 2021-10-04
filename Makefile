ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif



.PHONY: buildwin
buildwin: ## Run tests	
	GOOS=windows GOARCH=amd64 go build
