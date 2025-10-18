.PHONY: build dev package

build:
	wails3 build

dev: bindings
	wails3 dev

package:
	wails3 package

bindings:
	wails3 generate bindings