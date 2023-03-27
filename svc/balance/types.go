package wallet

import (
	"github.com/dmitrymomot/solana/metadata"
	"github.com/dmitrymomot/solana/types"
)

// Balance is a type that represents a token/SOL balance of a wallet.
// It includes the token metadata.
type Balance struct {
	Pubkey   string
	Mint     string
	IsNative bool
	Balance  types.TokenAmount
	Metadata *metadata.Metadata
}
