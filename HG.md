
### 文件系统
aufs
https://www.cnblogs.com/sammyliu/p/5931383.html
https://docs.docker.com/storage/storagedriver/select-storage-driver/

docker info 查看docker的storage

overlay2
https://www.jianshu.com/p/3826859a6d6e

### 使用busybox创建容器

docker export `docker run -itd busybox:latest` > busybox.tar
mkdir busybox && tar -xvf busybox.tar -C busybox



### /proc/self/exe

readlink -f /proc/$$/exe


https://stackoverflow.com/questions/606041/how-do-i-get-the-path-of-a-process-in-unix-linux



## Docker

	github.com/docker/docker v0.0.0-20171023200535-7848b8beb9d3
