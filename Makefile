.PHONY: build
build:
	@go build -o bin/ ./...

.PHONY: build-linux
build-linux:
	@GOOS=linux GOARCH=amd64 go build -o bin/ ./...

.PHONY: build-image
build-image: build-linux
	@docker-compose -f docker/docker-compose.yaml build

.PHONY: push
push: build-image
	@docker-compose -f docker/docker-compose.yaml push

.PHONY: install
install:
	@kubectl apply -k deploy

.PHONY: clean
clean:
	@kubectl delete -k deploy

.PHONY: install-test
install-test:
	@kubectl apply -f test/nginx.yaml

.PHONY: clean-test
clean-test:
	@kubectl delete -f test/nginx.yaml
