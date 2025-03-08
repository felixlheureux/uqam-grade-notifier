# Variables
APP_NAME="gnotifier"
DIST_FOLDER="../dist"

default: build

build:
	make clean
	go build -o $(DIST_FOLDER)/$(APP_NAME) main.go
	cp config.json $(DIST_FOLDER)
	cp grades.json $(DIST_FOLDER)

clean:
	rm -rf $(DIST_FOLDER)

.PHONY: default build clean
