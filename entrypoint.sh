#!/bin/sh
set -e

if [ "$1" = "server" ]; then
  echo "Running server...."
  ./server
  exit 0;
fi

if [ "$1" = "migrate" ]; then
  echo "Running migration...."
  ./migrate
  exit 0;
fi

echo "Running migration and server...."
./migrate
./server
exit 0