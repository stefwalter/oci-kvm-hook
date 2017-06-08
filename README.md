I'd like to have /dev/kvm available inside of each container started on
a host. Often kubernetes starts these containers. Once this works well
I'll submit an pull request the kubelet to make this a first class
concept there.

## Testing this out

Here are the manual commands to test this out. You should see a Linux boot process.

    $ sudo docker run -ti fedora /bin/bash
    # yum install wget qemu-kvm
    ...
    # wget https://download.fedoraproject.org/pub/alt/atomic/stable/Fedora-Atomic-25-20170314.0/CloudImages/x86_64/images/Fedora-Atomic-25-20170314.0.x86_64.qcow2
    # wget https://rawgit.com/cockpit-project/cockpit/master/test/common/cloud-init.iso
    # mknod /dev/kvm c 10 232
    # chmod 666 /dev/kvm
    # qemu-kvm -boot c -net nic -net user -m 1024 -nographic -cdrom cloud-init.iso Fedora-Atomic-25-20170314.0.x86_64.qcow2

Or with a prebuilt constainer:

    $ sudo docker run -ti stefwalter/test-kvm

Or in Openshift:

    $ oc create -f test/pod.json
    $ oc get pod test-kvm
    $ oc log test-kvm
