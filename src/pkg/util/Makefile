include $(GOROOT)/src/Make.inc

TARG=bugnuts/util
GOFILES=\
	timing.go\
	permute.go\
	math.go

include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

