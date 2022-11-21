#!/usr/bin/env bash

docker run \
  --net epgxml_net \
  --name epgxml \
  --rm -it \
  -v ${PWD}/tests/testdata/a4on.fdb:/firebird/data/a4on.fdb \
  -v ${PWD}/tests/testdata/TV_Pack.xml:/TV_Pack.xml \
  -v ${PWD}/config.yml:/config.yml \
  epgxml:latest
