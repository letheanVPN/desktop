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
	XMRig *XMRigMiner
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

// XMRigMiner represents an XMRig miner
type XMRigMiner struct {
	Miner
	ConfigPath string `json:"configPath"`
	API        *API   `json:"api"`
}

// API represents the XMRig API configuration
type API struct {
	Enabled    bool   `json:"enabled"`
	ListenHost string `json:"listenHost"`
	ListenPort int    `json:"listenPort"`
}

// XMRigSummary represents the summary from the XMRig API
type XMRigSummary struct {
	Hashrate struct {
		Total []float64 `json:"total"`
	} `json:"hashrate"`
	Results struct {
		SharesGood uint64 `json:"shares_good"`
		SharesTotal uint64 `json:"shares_total"`
	} `json:"results"`
	Uptime uint64 `json:"uptime"`
}
