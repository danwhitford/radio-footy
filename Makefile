GO_FILES := $(shell find . -name "*.go")

.PHONY: website
website: docs/index.html docs/icalendar.ics docs/styles.css docs/icon.png docs/log.txt docs/script.js

.PHONY: clean
clean:
	git clean -fdx

docs:
	mkdir -p docs

radiofooty: $(GO_FILES) cmd/radiofooty/website/template.go.tmpl cmd/radiofooty/website/icalendar.go.tmpl
	go build cmd/radiofooty/radiofooty.go

index.html icalendar.ics log.txt: radiofooty
	./radiofooty 2> log.txt

docs/index.html: docs index.html
	cp index.html docs/index.html

docs/log.txt: docs log.txt
	cp log.txt docs/log.txt

docs/icalendar.ics: docs icalendar.ics
	cp icalendar.ics docs/icalendar.ics

docs/%: static/%
	cp $< $@
