include $(GOROOT)/src/Make.inc

TARG=bugnuts/util
GOFILES=\
	gcdlcm.go\
	math.go\
	permute.go\

include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

