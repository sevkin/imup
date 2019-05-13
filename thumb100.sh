#!/bin/sh

# $1 input file
# $2 output file

convert "$1" -thumbnail 100x100^ -gravity center -extent 100x100 "$2"