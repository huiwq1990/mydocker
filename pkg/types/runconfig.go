package types

import "encoding/json"

type RunCommandConfig struct {
	Name string
	Volume string
	ImageName string
	Net string
	// 在容器中执行的命令，镜像的文件系统里面需要又这个可执行程序
	CmdArray []string
	Portmapping []string
	Tty bool
	Detach bool
	// cgroup相关
	Memory string
	Cpuset string
	Cpushare string
}

func (conf RunCommandConfig) String() string {
	bs,_ := json.Marshal(conf)
	return string(bs)
}