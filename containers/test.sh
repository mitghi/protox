#!/bin/bash

go test -v . && go test -run=^$ -bench=. -cpuprofile=cpu.out
