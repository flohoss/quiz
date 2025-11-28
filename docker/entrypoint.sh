#!/bin/sh
set -e

cat /app/logo.txt

mkdir -p /app/config
if [ -f /logo.svg ] && [ ! -f /app/config/logo.svg ]; then
	cp /logo.svg /app/config/logo.svg
fi
if [ -f /logo.svg ] && [ ! -f /app/config/favicon.svg ]; then
	cp /logo.svg /app/config/favicon.svg
fi

exec /app/quiz
