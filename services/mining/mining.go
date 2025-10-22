package mining

// Service is the mining service
//type Service struct {
//	runtime *wails.Runtime
//	log     *wails.CustomLogger
//	// miner         *Miner
//	// config        *Config
//	// stats         *Stats
//	// history       *History
//}

type Service struct {
	//config        *config.Config
	//networks      map[string]Network
	//activeNetwork Network
}

// Miner represents a miner
type Miner struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	URL           string `json:"url"`
	Path          string `json:"path"`
	Config        string `json:"config"`
	Running       bool   `json:"running"`
	Pid           int    `json:"pid"`
	LastHeartbeat int64  `json:"lastHeartbeat"`
}

// Config represents the config for a miner
type Config struct {
	Miner   string `json:"miner"`
	Pool    string `json:"pool"`
	Wallet  string `json:"wallet"`
	Threads int    `json:"threads"`
}

// Stats represents the stats for a miner
type Stats struct {
	Hashrate  int   `json:"hashrate"`
	Shares    int   `json:"shares"`
	Rejected  int   `json:"rejected"`
	Uptime    int   `json:"uptime"`
	LastShare int64 `json:"lastShare"`
}

// History represents the history of a miner
type History struct {
	Miner   string  `json:"miner"`
	Stats   []Stats `json:"stats"`
	Updated int64   `json:"updated"`
}
