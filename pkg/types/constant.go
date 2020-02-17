package types


const (
	ContainerRunning              = "running"
	ContainerStop                 = "stopped"
	ContainerExit                 = "exited"

	ContainerInfoLocation  = "/var/run/mydocker/"
	ContainerConfigName           = "config.json"
	ContainerLogFileName    string = "container.log"

	RootUrl				string = "/root"
	MntUrl				string = "/root/mnt/%s"
	ImageRepository  = "/root"

	WriteLayerUrl	string = "/root/writeLayer/"
)

