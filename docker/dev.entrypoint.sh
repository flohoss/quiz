#!/bin/sh
set -e

cat /logo.txt

air --build.cmd "go build -o tmp/bin/main ." \
	--build.include_ext "go" \
	--build.exclude_dir "web,docs,node_modules,dist,tmp" \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true
