
https://xionchen.github.io/2016/08/25/linux-bind-mount/

mkdir /tmp/nicktming
cd /tmp/nicktming
mkdir bin
which bash

# chroot需要路径下有/bin/bash
cp /bin/bash bin
cp /lib* .

wget https://busybox.net/downloads/binaries/1.21.1/busybox-x86_64

mv busybox-x86_64 busybox

chmod +x busybox



pwd

chroot /tmp/nicktming /bin/bash
./busybox ls -l