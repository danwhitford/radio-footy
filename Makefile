
.PHONY: cleancache
cleancache:
	find .cache/ -type f -delete

website: docs/index.html docs/icalendar.ics docs/lu.ttf
.PHONY: website

docs:
	mkdir -p docs

.cache:
	mkdir -p .cache

index.html: .cache
	go run cmd/radiofooty/radiofooty.go

docs/index.html: docs index.html
	cp index.html docs/index.html

icalendar.ics:
	go run cmd/radiofooty/radiofooty.go

docs/icalendar.ics: docs icalendar.ics
	cp icalendar.ics docs/icalendar.ics

docs/lu.ttf: docs
	cp lu.ttf docs/lu.ttf
