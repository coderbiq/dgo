#!/bin/bash

# god3 framework
mockgen -destination model_mocks.go \
   -package mocks \
   github.com/coderbiq/dgo/model \
   EventPublisher