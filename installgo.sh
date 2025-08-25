#!/bin/bash
Version=1.25.0
go install golang.org/dl/go$Version@latest
go$Version download
