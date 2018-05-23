GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: deps data build

data:
				mkdir -p ~/.local/m6ik/assets
				cp assets/*.json ~/.local/m6ik/assets/

build: 
				$(GOBUILD) -v ./cmd/m6ikgen
				$(GOBUILD) -v ./cmd/m6ikserv
test: 
				$(GOTEST) -v ./...
clean: 
				$(GOCLEAN)
				rm -f $(BINARY_NAME)
run:
				$(GOBUILD) -v ./cmd/m6ikgen
				./m6ikgen
deps:
				$(GOGET) github.com/kniren/gota/dataframe
				$(GOGET) github.com/kniren/gota/series
install:
				mv m6ikgen /usr/local/bin/
				mv m6ikserv /usr/local/bin/
