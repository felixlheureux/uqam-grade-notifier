# Variables
APP_NAME="gnotifier"
DIST_FOLDER="../dist"

default: build

build:
	make clean
	go build -o $(DIST_FOLDER)/$(APP_NAME) main.go

clean:
	rm -rf $(DIST_FOLDER)

.PHONY: default build clean
