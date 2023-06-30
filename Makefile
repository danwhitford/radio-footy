website: docs/index.html docs/icalendar.ics docs/lu.ttf
.PHONY: website

docs:
	mkdir -p docs

.cache:
	mkdir -p .cache

radiofooty:
	go build cmd/radiofooty/radiofooty.go

index.html: .cache radiofooty
	./radiofooty

docs/index.html: docs index.html
	cp index.html docs/index.html

icalendar.ics: .cache radiofooty
	./radiofooty

docs/icalendar.ics: docs icalendar.ics
	cp icalendar.ics docs/icalendar.ics

docs/lu.ttf: docs lu.ttf
	cp lu.ttf docs/lu.ttf
