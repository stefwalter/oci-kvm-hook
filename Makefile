.PHONY:  build install clean all docs dist

all: build docs

PREFIX ?= $(DESTDIR)/usr
HOOKSDIR=/usr/libexec/oci/hooks.d
HOOKSINSTALLDIR=$(DESTDIR)$(HOOKSDIR)
JSONDIR=/usr/share/containers/oci/hooks.d/
JSONINSTALLDIR=$(DESTDIR)$(JSONDIR)

# need this substitution to get build ID note
GOBUILD=go build -a -ldflags "${LDFLAGS:-} -B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')"

oci-kvm-hook: oci-kvm-hook.go
	GOPATH=$${GOPATH:+$$GOPATH:}/usr/share/gocode $(GOBUILD) -o oci-kvm-hook

oci-kvm-hook.1: oci-kvm-hook.1.md
	go-md2man -in "oci-kvm-hook.1.md" -out "oci-kvm-hook.1"
	sed -i 's|$$HOOKSDIR|$(HOOKSDIR)|' oci-kvm-hook.1

docs: oci-kvm-hook.1
build: oci-kvm-hook

# Must be run with: sudo make install
install: oci-kvm-hook oci-kvm-hook.1
	install -d -m 755 $(HOOKSINSTALLDIR)
	install -m 755 oci-kvm-hook $(HOOKSINSTALLDIR)
	install -m 755 oci-kvm-hook.json $(JSONINSTALLDIR)
	install -d -m 755 $(PREFIX)/share/man/man1
	install -m 644 oci-kvm-hook.1 $(PREFIX)/share/man/man1

# Must be run with: sudo make test
test: install
	test/smoke-test.sh

clean:
	rm -f oci-kvm-hook *~
	rm -f oci-kvm-hook.1
	rm -f *.tar.gz *.rpm
	rm -rf ./x86_64/

rpms:
	sh -c "git archive HEAD --prefix=oci-kvm-hook-$$(git describe --abbrev=0)/ --output=$$(git describe --abbrev=0).tar.gz"
	rpmbuild -ba --define "_sourcedir $$PWD" --define "_specdir $$PWD" --define "_rpmdir $$PWD" --define "_srcrpmdir $$PWD" oci-kvm-hook.spec
