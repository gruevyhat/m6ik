GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=m6ik

all: deps data build

data:
				mkdir -p ~/.local/m6ik/assets
				cp assets/*.json ~/.m6ik/assets

build: 
				$(GOBUILD) -v -o $(BINARY_NAME) ./cmd/chargen
test: 
				$(GOTEST) -v ./...
clean: 
				$(GOCLEAN)
				rm -f $(BINARY_NAME)
run:
				$(GOBUILD) -o $(BINARY_NAME) ./cmd/chargen
				./$(BINARY_NAME)
deps:
				$(GOGET) github.com/kniren/gota/dataframe
				$(GOGET) github.com/kniren/gota/series
