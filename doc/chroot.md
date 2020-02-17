

pivot_root和chroot的区别：

pivot_root把进程切换到一个新的root目录，对之前root文件系统的不再有依赖，这样你就能够umount原先的root文件系统。
chroot只是更改了root目录，还会依赖老的文件系统,主要使用的目的一般是为了限制用户的访问。


chroot MyDir, 命令执行后，默认会执行${SHELL} -i。简单的rootfilesystem没有bash命令，可以执行 chroot MyDir /bin/sh


https://www.ibm.com/developerworks/cn/linux/l-cn-chroot/index.html
https://juejin.im/post/5c2b495af265da6134388142