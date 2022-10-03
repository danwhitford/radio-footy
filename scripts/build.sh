#!/bin/bash
set -ux

mkdir -p docs
cp index.html docs/
cp icalendar.ics docs/
cp -r static/* docs/
