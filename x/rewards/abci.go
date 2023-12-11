package rewards

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lavanet/lava/x/rewards/keeper"
)

// BeginBlocker calculates the validators block rewards and transfers them to the fee collector
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.DistributeBlockReward(ctx)
	if err != nil {
		panic(err)
	}
}
