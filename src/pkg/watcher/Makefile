include $(GOROOT)/src/Make.inc

TARG=bugnuts/watcher
GOFILES=\
	watcher.go

include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

