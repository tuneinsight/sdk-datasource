#!/usr/bin/env bash

VERSION=$(git describe --tags 2> /dev/null)
if [ $? -eq 0 ]; then
  echo "$VERSION"
else
  echo -n 0.0.1-
  echo -n $(git branch --show-current)-
  echo $(git describe --tags --always)
fi
