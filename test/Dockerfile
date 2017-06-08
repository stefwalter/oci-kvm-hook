FROM fedora:25

RUN curl -o /atomic.qcow2 https://ftp-stud.hs-esslingen.de/pub/Mirrors/alt.fedoraproject.org/atomic/stable/Fedora-Atomic-26-20170707.1/CloudImages/x86_64/images/Fedora-Atomic-26-20170707.1.x86_64.qcow2
RUN yum -y install qemu-kvm
ADD cloud-init.iso /
CMD test -e /dev/kvm || mknod --mode=0666 /dev/kvm c 10 232; qemu-kvm -boot c -net nic -net user -m 1024 -nographic -cdrom /cloud-init.iso /atomic.qcow2
