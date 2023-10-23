package keeper_test

import (
	"github.com/stretchr/testify/suite"
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/testing/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx      sdk.Context
	msKeeper *multistakingkeeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.ctx, suite.msKeeper = ctx, &app.MultiStakingKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
