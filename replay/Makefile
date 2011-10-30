include $(GOROOT)/src/Make.inc

TARG=replay
GOFILES=\
	replay.go


include $(GOROOT)/src/Make.pkg


.PHONY: gofmt test
gofmt:
	gofmt -w $(GOFILES)

