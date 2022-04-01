#!/bin/bash
set -ux

mkdir -p docs
cp index.html docs/
cp -r static/* docs/
