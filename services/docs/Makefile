.PHONY: serve build clean collect-docs

serve:
	mkdocs serve

build:
	mkdocs build -d static

collect-docs:
	mkdir -p docs static
	cp -r src/docs/assets docs/
	cp -r src/docs/getting-started docs/
	cp -r src/docs/sdk docs/
	cp -r src/docs/updates docs/

clean:
	rm -rf dist docs
