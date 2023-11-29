package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type Keeper struct {
	storeKey         storetypes.StoreKey
	memKey           storetypes.StoreKey
	cdc              codec.BinaryCodec
	stakingKeeper    types.StakingKeeper
	stakingMsgServer stakingtypes.MsgServer
	distrMsgServer   distrtypes.MsgServer
	bankKeeper       types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	stakingKeeper stakingkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	bankKeeper types.BankKeeper,
	key storetypes.StoreKey,
	memKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		storeKey:         key,
		memKey:           memKey,
		stakingKeeper:    stakingKeeper,
		stakingMsgServer: stakingkeeper.NewMsgServerImpl(stakingKeeper),
		distrMsgServer:   distrkeeper.NewMsgServerImpl(distrKeeper),
		bankKeeper:       bankKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
