# Variables
DIST_FOLDER="dist"

default: build

build: build-server build-cron

build-server:
	cd server && make build

build-cron:
	cd cron && make build

clean:
	rm -rf $(DIST_FOLDER)
	cd server && make clean
	cd cron && make clean

run-server:
	cd server && make run

run-cron:
	cd cron && make run

.PHONY: default build build-server build-cron clean run-server run-cron 