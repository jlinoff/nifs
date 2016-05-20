#!/bin/bash
#
# Run the program, see if it accepts all options correctcly.
#
go build nifs.go
set -ex
./nifs -h
./nifs --help
./nifs -V
./nifs --version

./nifs
./nifs -4
./nifs --ipv4
./nifs -6
./nifs --ipv4
./nifs -u
./nifs --unicast
./nifs -m
./nifs --multicast
./nifs -H
./nifs --hw
./nifs -s
./nifs --sw

set +ex
echo "success - all switches recognized"
