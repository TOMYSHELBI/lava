package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/lavanet/lava/x/pairing/client/cli"
)

var PairingUnstakeProposal = govclient.NewProposalHandler(cli.NewSubmitUnstakeProposalTxCmd)
