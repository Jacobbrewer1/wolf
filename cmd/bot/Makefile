# Define variables
hash = $(shell git rev-parse --short HEAD)
registry = ghcr.io/jacobbrewer1/wolf

# Define targets
all: linux windows macos

CI:
	@echo "Building for production"
	make deps
	GOOS=linux GOARCH=amd64 go build -o bin/app -tags prod
	docker build -t $(registry):$(hash) .
	docker push $(registry):$(hash)
	# Clean up
	make clean
linux:
	@echo "Building for linux"
	GOOS=linux GOARCH=amd64 go build -o bin/app
windows:
	@echo "Building for windows"
	GOOS=windows GOARCH=amd64 go build -o bin/app.exe
macos:
	@echo "Building for mac os"
	GOOS=darwin GOARCH=amd64 go build -o bin/app_mac
clean:
	@echo "Cleaning up"
	rm -rf bin
deps:
	@echo "Getting dependencies"
	go mod tidy
	go mod vendor
	go mod verify
# The docker command takes the build tag as an argument
docker:
	@echo "Building docker image"
	make linux
	# Build the docker image
	docker build -t $(registry):$(hash) .
	docker scout quickview $(registry):$(hash)
	docker scout cves $(registry):$(hash)
	docker scout recommendations $(registry):$(hash)
	# Clean up
	make clean
