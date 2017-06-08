FROM fedora:25

RUN curl -o /image.qcow2 https://dl.fedoraproject.org/pub/fedora/linux/releases/25/CloudImages/x86_64/images/Fedora-Cloud-Base-25-1.3.x86_64.qcow2
RUN yum -y install qemu-kvm
ADD cloud-init.iso /
CMD test -f /dev/kvm || mknod --mode=0666 /dev/kvm c 10 232; qemu-kvm -boot c -net nic -net user -m 1024 -nographic -cdrom /cloud-init.iso /image.qcow2
