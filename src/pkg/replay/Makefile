include $(GOROOT)/src/Make.inc

TARG=bugnuts/replay
GOFILES=\
	load.go\
	metadata.go\
	replay.go\


include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

