https://www.jianshu.com/p/3826859a6d6e
https://blog.csdn.net/zhonglinzhang/article/details/80970411
https://blog.csdn.net/luckyapple1028/article/details/78075358


mkdir -p /tmp/testvolume
touch -p /tmp/testvolume/fromvolume

mkdir -p /tmp/overlay2_test
cd /tmp/overlay2_test
mkdir -p lower1 lower2 upper worker merge
mkdir -p upper/app
mkdir -p lower1/dir lower2/dir upper/dir
touch lower1/foo1 lower2/foo2 upper/foo3
touch lower1/dir/aa lower2/dir/aa lower1/dir/bb upper/dir/bb
echo "from lower1" >> lower1/dir/aa
echo "from lower2" >> lower2/dir/aa
echo "from lower1" >> lower1/dir/bb
echo "from upper" >> upper/dir/bb

mount --bind /root/ /tmp/overlay2_test/upper/app

mount -t overlay overlay -o lowerdir=lower1:lower2,upperdir=upper,workdir=worker merge


# dir foo1 foo2 foo3
ls merge
# aa bb
ls merge/dir

# from lower1
cat merge/dir/aa

# from upper
cat upper/dir/bb



## 测试chroot能否成功
```
tmpPath=/root/writeLayer/xyz
mkdir -p ${tmpPath}
cd ${tmpPath}
mkdir -p upper worker merge
mount -t overlay overlay -o lowerdir=/root/busybox,upperdir=${tmpPath}/upper,workdir=${tmpPath}/worker ${tmpPath}/merge
cd merge
mkdir oldroot
unshare -m
pivot_root ./ oldroot
```


```
mkdir /new-root
# -t：指定档案系统的型态
mount -n -t tmpfs -o size=500M none /new-root
cd /new-root
mkdir old-root
#解决pivot_root命令无效的问题
unshare -m
pivot_root . old-root
```


参考
https://www.cnblogs.com/bianhao3321/p/6873511.html
