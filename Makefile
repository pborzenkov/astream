GO ?= go

all: astream

astream: FORCE
	go install ./cmd/astream

FORCE:
