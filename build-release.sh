#!/bin/bash
if [ ! -d "build" ]; then
  mkdir build
fi
pv="0.0.3"
pn="obs-drops-overlay"
GOARCH="amd64"
declare -a goos=("windows" "linux" "darwin")
for i in "${goos[@]}"; do
  if [ "${i}" == "windows" ]; then
    CGO_ENABLED=0 GOOS="${i}" go build -v -a -ldflags="-s -w" -o "build/${pn}-${pv}-${i}-${GOARCH}.exe" .
  else
    CGO_ENABLED=0 GOOS="${i}" go build -v -a -ldflags="-s -w" -o "build/${pn}-${pv}-${i}-${GOARCH}" .
  fi
done
for f in build/*; do
  file "${f}"
done
