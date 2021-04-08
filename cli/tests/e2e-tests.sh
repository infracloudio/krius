#!/bin/sh

kind create cluster --name krius

clean_up() {
  ec=$?
  kind delete cluster --name krius
  exit $ec
}

trap clean_up SIGHUP SIGINT SIGTERM
go test -v ./tests/krius-cli-tests --timeout 45m
clean_up