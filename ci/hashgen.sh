#!/bin/sh

for f in inlets*; do shasum -a 256 $f > $f.sha256; done
