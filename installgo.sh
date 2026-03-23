#!/bin/bash
Version=1.26.1
go install golang.org/dl/go$Version@latest
go$Version download
