.PHONY: build dev package test

build:
	wails3 build

dev: bindings
	wails3 dev

package:
	wails3 package

test:
	go test ./services/*

bindings:
	wails3 generate bindings -ts