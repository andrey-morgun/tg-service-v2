package app

import (
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/one/chain_client"
	"log"
)

func (a *App) initChainClient() {
	a.chain = chain_client.NewClient(a.config.Chain)
	if a.chain == nil {
		log.Fatal(errs.InternalError{Cause: "chain-client = nil"})
	}
}
