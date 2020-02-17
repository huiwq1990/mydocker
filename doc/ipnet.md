


```
#帮助命令
ip netns help

#查看ns
ip netns list
ls -l /var/run/netns

#增加一个network namespace: neta
ip netns add neta

#进入到neta namespace
ip netns exec neta sh
# ifconfig

#查看某个pid的namespace
ip netns identify 28033

#查看neta namespace下的所有进程pid
ip netns pids neta

ip netns monitor

ip netns delete neta
ls -l /var/run/netns
ip netns list
```