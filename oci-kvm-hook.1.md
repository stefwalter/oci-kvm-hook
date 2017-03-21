% OCI-KVM-HOOK(1) oci-kvm-hook
% November 2015
## NAME
oci-kvm-hook - Make /dev/kvm available in any OCI container

## SYNOPSIS

**oci-kvm-hook**

## DESCRIPTION

`oci-kvm-hook` is a OCI hook program. If you add it to the runc json data
as a hook, runc will execute the application after the container process
is created but before it is executed, with a `prestart` flag.  When the
container starts `oci-kvm-hook` will make sure that `/dev/kvm` is available
inside the container.

## EXAMPLES

	$ docker run -it busybox /bin/sh
	$ ls -l /dev/kvm

## SEE ALSO

docker-run(1)
