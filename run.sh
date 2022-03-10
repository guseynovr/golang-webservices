#!/bin/bash

src=/Users/dgidget/Documents/go/src/coursera/coursera_hw
task=${1:-hw4_test_coverage}
image=${2:-golang:1.9.2}
docker run -v $src/$task:/go -it --rm --name mailgo $image

