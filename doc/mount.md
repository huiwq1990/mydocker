
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


mkdir /run/something
cd /run/something
mkdir -p etc/something lib usr/lib usr/sbin var/lib/something bin
mount --bind /lib lib
mount --bind /usr/lib usr/lib
mount --bind /usr/sbin usr/sbin
mount --bind /bin bin

mount -o remount,ro,bind lib
mount -o remount,ro,bind usr/lib
mount -o remount,ro,bind usr/sbin
mount -o remount,ro,bind bin

mkdir -p /app
touch /app/aaaa
mkdir -p xxx
mount --bind /app /root/bbb/xxx
mount -o remount,ro,bind /root/bbb/xxx

chroot . /bin/ls &