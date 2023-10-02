GO_FILES := $(shell find . -name "*.go")

.PHONY: website
website: docs/index.html docs/icalendar.ics docs/styles.css docs/icon.png

.PHONY: clean
clean:
	git clean -fdx

docs:
	mkdir -p docs

.cache:
	mkdir -p .cache

radiofooty: $(GO_FILES)
	go build cmd/radiofooty/radiofooty.go

index.html: .cache radiofooty internal/website/template.go.tmpl
	./radiofooty

docs/index.html: docs index.html
	mv index.html docs/index.html

icalendar.ics: .cache radiofooty internal/website/icalendar.go.tmpl
	./radiofooty

docs/icalendar.ics: docs icalendar.ics
	mv icalendar.ics docs/icalendar.ics

docs/%: static/%
	cp $< $@
