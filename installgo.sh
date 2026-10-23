#!/bin/bash
Version=1.27.3
go install golang.org/dl/go$Version@latest
go$Version download
