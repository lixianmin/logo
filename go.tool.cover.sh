#!/bin/bash

FILENAME=coverage.out
go test -coverprofile=$FILENAME
go tool cover -html=$FILENAME
rm $FILENAME

