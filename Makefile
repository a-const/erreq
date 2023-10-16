PACKAGE_NAME = erreq
VERSION = 0.4
DESTDIR = ./build/pkg

.PHONY: clean

all: build package clean

build:
	go build -o ./build/erreq

package:
	mkdir -p $(DESTDIR)/usr/local/bin
	mkdir -p $(DESTDIR)/DEBIAN
	cp ./build/erreq $(DESTDIR)/usr/local/bin
	cp control $(DESTDIR)/DEBIAN
	dpkg-deb --build $(DESTDIR) $(PACKAGE_NAME)_$(VERSION)_amd64.deb

clean:
	rm -rf build/pkg
