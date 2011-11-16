include $(GOROOT)/src/Make.inc

TARG=bugnuts/debug
GOFILES=\
	debugging.go

include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

