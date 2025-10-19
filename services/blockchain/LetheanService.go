package lthn

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	lthn "github.com/letheanVPN/blockchain/utils/sdk/client/go"
)

// Brand defines the type for different application brands
type Brand string

const (
	AdminHub     Brand = "admin-hub"
	ServerHub    Brand = "server-hub"
	GatewayHub   Brand = "gateway-hub"
	DeveloperHub Brand = "developer-hub"
	ClientHub    Brand = "client-hub"
)

type LetheanService struct {
}

func NewLetheanService() *LetheanService {
	return &LetheanService{}
}

// FetchBlockData Fetch blockchain block data
func (s *LetheanService) FetchBlockData() {
	identifier := "1000" // string | The hash (hex string) or height (integer) of the block to retrieve.

	configuration := lthn.NewConfiguration()
	apiClient := lthn.NewAPIClient(configuration)
	resp, r, err := apiClient.BlockUtilsSdkClientGo.GetBlock(context.Background(), identifier).Execute()
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error when calling `BlockUtilsSdkClientGo.GetBlock``: %v\n", err)
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		if err != nil {
			return
		}
	}
	// response from `GetBlock`: BlockDetailsModel
	_, err = fmt.Fprintf(os.Stdout, "Response from `BlockUtilsSdkClientGo.GetBlock`: %v\n", resp)
	if err != nil {
		return
	}
}
