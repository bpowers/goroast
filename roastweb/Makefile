EXE:=$(shell basename `pwd`)
REMOTE:=rpi.local

all: $(EXE)

$(EXE):
	go test
	go build
	mkdir -p site
	cp $(EXE) site/
	rsync -av add-user err static tmpl site

put: $(EXE)
	rsync -cavz site/ $(REMOTE):site-$(EXE)

clean:
	rm -rf site $(EXE)

.PHONY: $(EXE) clean
