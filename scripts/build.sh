#!/bin/bash
set -ux

mkdir -p docs
cp index.html docs/
cp icalendar docs/
cp -r static/* docs/
