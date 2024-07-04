module go-tha-utxos

go 1.22

require (
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.19.0 // indirect
)

replace github.com/tendermint/tendermint => github.com/cometbft/cometbft v0.34.27
