.PHONY: clean
clean:
	git clean -fdx

website: docs/index.html docs/icalendar.ics docs/lu.ttf
.PHONY: website

docs:
	mkdir -p docs

.cache:
	mkdir -p .cache

radiofooty: $(wildcard cmd/radiofooty/*.go) $(wildcard internal/feeds/*.go) $(wildcard internal/filecacher/*.go) $(wildcard internal/website/*.tmpl)
	go build cmd/radiofooty/radiofooty.go

index.html: .cache radiofooty
	./radiofooty

docs/index.html: docs index.html
	mv index.html docs/index.html

icalendar.ics: .cache radiofooty
	./radiofooty

docs/icalendar.ics: docs icalendar.ics
	mv icalendar.ics docs/icalendar.ics

docs/lu.ttf: docs lu.ttf
	cp lu.ttf docs/lu.ttf
