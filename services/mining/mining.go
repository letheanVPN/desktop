package mining

import (
	"net/http"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
)

// Miner is the interface for a miner
type Miner interface {
	Install() error
	Start(config *Config) error
	Stop() error
	GetStats() (*PerformanceMetrics, error)
	GetName() string
}

type Service struct {
	Miners map[string]Miner
	Router *gin.Engine
	Server *http.Server
}

// Config represents the config for a miner
type Config struct {
	Miner     string `json:"miner"`
	Pool      string `json:"pool"`
	Wallet    string `json:"wallet"`
	Threads   int    `json:"threads"`
	TLS       bool   `json:"tls"`
	HugePages bool   `json:"hugePages"`
}

// PerformanceMetrics represents the performance metrics for a miner
type PerformanceMetrics struct {
	Hashrate  int                    `json:"hashrate"`
	Shares    int                    `json:"shares"`
	Rejected  int                    `json:"rejected"`
	Uptime    int                    `json:"uptime"`
	LastShare int64                  `json:"lastShare"`
	Algorithm string                 `json:"algorithm"`
	ExtraData map[string]interface{} `json:"extraData,omitempty"`
}

// History represents the history of a miner
type History struct {
	Miner   string               `json:"miner"`
	Stats   []PerformanceMetrics `json:"stats"`
	Updated int64                `json:"updated"`
}

// XMRigMiner represents an XMRig miner
type XMRigMiner struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	URL           string `json:"url"`
	Path          string `json:"path"`
	Running       bool   `json:"running"`
	LastHeartbeat int64  `json:"lastHeartbeat"`
	ConfigPath    string `json:"configPath"`
	API           *API   `json:"api"`
	mu            sync.Mutex
	cmd           *exec.Cmd `json:"-"`
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
		SharesGood  uint64 `json:"shares_good"`
		SharesTotal uint64 `json:"shares_total"`
	} `json:"results"`
	Uptime    uint64 `json:"uptime"`
	Algorithm string `json:"algo"`
}
