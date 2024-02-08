GO_FILES := $(shell find . -name "*.go")

.PHONY: website
website: docs/index.html docs/icalendar.ics docs/styles.css docs/icon.png docs/log.txt

.PHONY: clean
clean:
	git clean -fdx

docs:
	mkdir -p docs

radiofooty: $(GO_FILES)
	go build cmd/radiofooty/radiofooty.go

index.html log.txt: radiofooty internal/website/template.go.tmpl
	./radiofooty website > index.html 2> log.txt

docs/index.html: docs index.html
	mv index.html docs/index.html

docs/log.txt: docs log.txt
	mv log.txt docs/log.txt

icalendar.ics: radiofooty internal/website/icalendar.go.tmpl
	./radiofooty calendar > icalendar.ics

docs/icalendar.ics: docs icalendar.ics
	mv icalendar.ics docs/icalendar.ics

docs/%: static/%
	cp $< $@
