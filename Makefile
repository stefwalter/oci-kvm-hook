.PHONY:  build install clean all docs dist

all: build docs

PREFIX ?= $(DESTDIR)/usr
HOOKSDIR=/usr/libexec/oci/hooks.d
HOOKSINSTALLDIR=$(DESTDIR)$(HOOKSDIR)

# need this substitution to get build ID note
GOBUILD=go build -a -ldflags "${LDFLAGS:-} -B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')"

oci-kvm-hook: oci-kvm-hook.go
	GOPATH=$${GOPATH:+$$GOPATH:}/usr/share/gocode $(GOBUILD) -o oci-kvm-hook

oci-kvm-hook.1: oci-kvm-hook.1.md
	go-md2man -in "oci-kvm-hook.1.md" -out "oci-kvm-hook.1"
	sed -i 's|$$HOOKSDIR|$(HOOKSDIR)|' oci-kvm-hook.1

docs: oci-kvm-hook.1
build: oci-kvm-hook

install: oci-kvm-hook oci-kvm-hook.1
	install -d -m 755 $(HOOKSINSTALLDIR)
	install -m 755 oci-kvm-hook $(HOOKSINSTALLDIR)
	install -d -m 755 $(PREFIX)/share/man/man1
	install -m 644 oci-kvm-hook.1 $(PREFIX)/share/man/man1

clean:
	rm -f oci-kvm-hook *~
	rm -f oci-kvm-hook.1
