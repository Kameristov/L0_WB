LOCAL_DIR:=$(CURDIR)

run: # run app
	go run "$(LOCAL_DIR)/cmd/app/main.go"
.PHONY: run 

build: # build app
	go build -C $(LOCAL_DIR)/cmd/app -o $(LOCAL_DIR)/L0_WB.exe main.go
.PHONY: build 

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

mock: ### run mockgen
	mockgen -source ./internal/usecase/interfaces.go -package usecase_test > ./internal/usecase/mocks_test.go
.PHONY: mock