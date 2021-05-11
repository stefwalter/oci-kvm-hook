#!/bin/sh

set -eux

# This smoke test reqires:
#  1. The oci-kvm-hook is installed
#  2. It's run as root

podman run -ti --rm fedora ls -l /dev | tee /tmp/oci-kvm-dev.log
if ! grep -w 'kvm' /tmp/oci-kvm-dev.log; then
	echo "TEST FAILED"
	exit 1
fi
