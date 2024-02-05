package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/strangelove-ventures/balances/cosmos"
)

type ChainWallet struct {
	Chain         string
	RPCURL        string
	WalletAddress string
}

var wallets = []ChainWallet{
	{
		Chain:         "cosmoshub-4",
		RPCURL:        "https://rpc.cosmoshub.strange.love:443",
		WalletAddress: "cosmos130mdu9a0etmeuw52qfxk73pn0ga6gawkryh2z6",
	},
	{
		Chain:         "osmosis-1",
		RPCURL:        "https://rpc.osmosis.strange.love:443",
		WalletAddress: "osmo1r2u5q6t6w0wssrk6l66n3t2q3dw2uqny0l6fwk",
	},
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx := context.Background()

	for _, wallet := range wallets {
		cc, err := cosmos.NewProvider(wallet.RPCURL)
		if err != nil {
			log.Error("Failed to create cosmos provider", "error", err)
			continue
		}

		// query latest height
		height := uint64(0)

		coins, err := cc.QueryBalance(ctx, wallet.WalletAddress, height)
		if err != nil {
			log.Error(
				"Failed to query balance",
				"wallet", wallet.WalletAddress,
				"chain", wallet.Chain,
				"rpc_url", wallet.RPCURL,
				"error", err,
			)
		}

		for _, coin := range coins {
			log.Info(
				"Balance",
				"wallet", wallet.WalletAddress,
				"chain", wallet.Chain,
				"rpc_url", wallet.RPCURL,
				"coin", coin.String(),
			)
		}
	}
}
