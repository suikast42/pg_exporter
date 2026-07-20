#!/bin/bash
Version=1.26.5
go install golang.org/dl/go$Version@latest
go$Version download
