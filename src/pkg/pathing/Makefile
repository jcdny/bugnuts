include $(GOROOT)/src/Make.inc

TARG=bugnuts/pathing
GOFILES=\
	fill.go\
	pathing.go\
	queue.go\


include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

