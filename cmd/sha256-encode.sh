#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <string>"
  exit 1
fi

string="$1"

echo -n "$string" | sha256sum | awk '{print $1}'
