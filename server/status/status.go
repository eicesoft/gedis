package status

// ServerStatus 服务器状态
type ServerStatus struct {
	// Server 服务器信息
	Server struct {
		version         string
		mode            string
		os              string
		archBits        int
		processId       int
		tcpPort         string
		uptimeInSeconds int
		uptimeInDays    int
		configFile      string
		executable      string
	}

	//
	Clients struct {
		connectedClients int
		trackingClients  int
	}

	Memory struct {
		usedMemory int32
	}

	Stats struct {
		totalConnectionsReceived int64
		totalCommandsProcessed   int64
	}
}
