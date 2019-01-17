#!/bin/bash

# god3 framework
mockgen -destination devent_mocks.go \
   -package mocks \
   github.com/coderbiq/dgo/base/devent \
   Publisher