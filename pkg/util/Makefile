include $(GOROOT)/src/Make.inc

TARG=bugnuts/util
GOFILES=\
	timing.go\
	util.go

include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

