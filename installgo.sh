#!/bin/bash
Version=1.26.2
go install golang.org/dl/go$Version@latest
go$Version download
