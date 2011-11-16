include $(GOROOT)/src/Make.inc

TARG=bugnuts/maps
GOFILES=\
	data.go\
	direction.go\
	item.go\
	map.go\
	mask.go\
	point.go\
	torus.go\
	sym.go\
	util.go


include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

