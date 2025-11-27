#!/bin/sh
set -e

cat /app/logo.txt

if [ -f /logo.svg ] && [ ! -f /app/config/logo.svg ]; then
	cp /logo.svg /app/config/logo.svg
fi

exec /app/quiz
