package action

import (
	"fmt"
	"go-tha-utxos/app/services"
	"go-tha-utxos/config"
)

func RunApp() error {
	cfg, err := config.LoadAndApplyConfig()
	if err != nil {
		return fmt.Errorf("could not load config: %v", err)
	}

	daemon, err := services.NewRpcDaemon(&cfg.RpcConnection)
	if err != nil {
		return fmt.Errorf("could not create rpc daemon: %v", err)
	}

	addresses, err := daemon.GetNewAddresses(1)
	if err != nil {
		return fmt.Errorf("could not get addresses: %v", err)
	}

	_ = addresses

	return nil
}
