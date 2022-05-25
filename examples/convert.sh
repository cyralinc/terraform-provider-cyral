#!/bin/bash

FILES="./resources/*.md"
for f in $FILES
do
  #echo "Processing $f file..."
  filename=$(basename "$f")
  fname="${filename%.*}"
  echo $fname
  dirname=./resources/$fname
  mkdir -p $dirname
  mv $f $dirname/resource.tf
  # take action on each file. $f store current file name
  # perform some operation with the file
done
