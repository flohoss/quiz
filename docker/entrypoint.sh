#!/bin/sh
set -e

cat /app/mono-12.txt

mkdir -p /app/config
if [ -f /logo.svg ] && [ ! -f /app/config/logo.svg ]; then
	cp /logo.svg /app/config/logo.svg
fi

exec /app/quiz
