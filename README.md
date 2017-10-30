I'd like to have /dev/kvm available inside of each container started on
a host. Often kubernetes starts these containers. Once this works well
I'll submit an pull request the kubelet to make this a first class
concept there.

## Testing this out

Here are the manual commands to test this out. You should see a Linux boot process.

    $ sudo docker run -ti fedora /bin/bash
    # yum install qemu-kvm
    ...
    # curl -Lo atomic.qcow2 https://ftp-stud.hs-esslingen.de/pub/Mirrors/alt.fedoraproject.org/atomic/stable/Fedora-Atomic-26-20170707.1/CloudImages/x86_64/images/Fedora-Atomic-26-20170707.1.x86_64.qcow2
    # curl -Lo cloud-init.iso https://rawgit.com/stefwalter/oci-kvm-hook/master/test/cloud-init.iso
    # qemu-kvm -boot c -net nic -net user -m 1024 -nographic -cdrom cloud-init.iso atomic.qcow2

Or with a prebuilt constainer:

    $ sudo docker run -ti --rm stefwalter/test-kvm

Or in Openshift:

    $ oc create -f test/pod.json
    $ oc get pod test-kvm
    $ oc log test-kvm

(This assumes the pod will run in an SCC with RunAsAny).
