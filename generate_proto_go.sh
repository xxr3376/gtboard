#!/bin/bash -e

files=$(find proto -name "*.proto" | sed 's/proto\///g' | sed 's/\.proto//g')

for name in $files; do
    file="proto/$name.proto"
    protoc --proto_path=proto --go_opt="module=github.com/xxr3376/gtboard" --go_out=. $file
done
