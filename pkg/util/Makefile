include $(GOROOT)/src/Make.inc

TARG=bugnuts/util
GOFILES=\
	util.go

include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

