package app

import (
	"github.com/andReyM228/one/chain_client"
)

func (a *App) initChainClient() {
	a.chain = chain_client.NewClient(a.config.Chain)
}
