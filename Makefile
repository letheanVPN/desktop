.PHONY: build dev package test

build:
	wails3 build

dev: bindings
	wails3 dev

package:
	wails3 package

test:
	go test ./services/*

fetch-reviews:
ifeq ($(PR),)
	$(error PR is not set. Usage: make fetch-reviews PR=<number>)
endif
	gh pr view $(PR) --json reviews > reviews.json
	task

bindings:
	wails3 generate bindings -ts