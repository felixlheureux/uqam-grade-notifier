# Variables
APP_NAME="gnotifier-server"
DIST_FOLDER="../dist"

default: build

build:
	make clean
	go build -o $(DIST_FOLDER)/$(APP_NAME) main.go

clean:
	rm -rf $(DIST_FOLDER)

run:
	./$(DIST_FOLDER)/$(APP_NAME) -config config.json

.PHONY: default build clean run 