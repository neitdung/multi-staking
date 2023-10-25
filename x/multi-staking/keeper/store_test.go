package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/testutil"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
)

func (suite *KeeperTestSuite) TestSetBondTokenWeight() {
	suite.SetupTest()

	gasDenom := "ario"
	govDenom := "arst"
	gasWeight := math.LegacyNewDec(1)
	govWeight := math.LegacyNewDecWithPrec(2, 4)

	suite.msKeeper.SetBondTokenWeight(suite.ctx, gasDenom, gasWeight)
	suite.msKeeper.SetBondTokenWeight(suite.ctx, govDenom, govWeight)

	suite.Equal(gasWeight, suite.msKeeper.GetBondTokenWeight(suite.ctx, gasDenom))
	suite.Equal(govWeight, suite.msKeeper.GetBondTokenWeight(suite.ctx, govDenom))
}

func (suite *KeeperTestSuite) TestSetValidatorBondDenom() {
	valA := testutil.GenValAddress()
	valB := testutil.GenValAddress()
	gasDenom := "ario"
	govDenom := "arst"
	testCases := []struct {
		name string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string
		vals []sdk.ValAddress
		expPanic bool
	}{
		{
			name: "1 val, 1 denom, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorBondDenom(ctx, valA, gasDenom)
				return []string{gasDenom}
			},
			vals: []sdk.ValAddress{valA},
			expPanic: false,
		},
		{
			name: "2 val, 2 denom, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorBondDenom(ctx, valA, gasDenom)
				msKeeper.SetValidatorBondDenom(ctx, valB, govDenom)
				return []string{gasDenom, govDenom}
			},
			vals: []sdk.ValAddress{valA, valB},
			expPanic: false,
		},
		{
			name: "1 val, 2 denom, failed",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorBondDenom(ctx, valA, gasDenom)
				msKeeper.SetValidatorBondDenom(ctx, valA, govDenom)
				return []string{gasDenom, govDenom}
			},
			vals: []sdk.ValAddress{valA, valB},
			expPanic: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
		
			if tc.expPanic {
				suite.Require().PanicsWithValue("validator denom already set",func() {
					tc.malleate(suite.ctx, suite.msKeeper)
				})
			} else {
				inputs := tc.malleate(suite.ctx, suite.msKeeper)
				for idx, val := range tc.vals {
					actualDenom := suite.msKeeper.GetValidatorBondDenom(suite.ctx, val)
					suite.Require().Equal(inputs[idx], actualDenom)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSetIntermediaryAccountDelegator() {
	delA := testutil.GenAddress()
	delB := testutil.GenAddress()
	imAddrressA := testutil.GenAddress()
	imAddrressB := testutil.GenAddress()

	testCases := []struct {
		name string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.AccAddress
		imAccs []sdk.AccAddress
		expPanic bool
	}{
		{
			name: "1 delegator, 1 intermediary account, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.AccAddress {
				msKeeper.SetIntermediaryAccountDelegator(ctx, imAddrressA, delA)
				return []sdk.AccAddress{delA}
			},
			imAccs: []sdk.AccAddress{imAddrressA},
			expPanic: false,
		},
		{
			name: "2 delegator, 2 intermediary account, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.AccAddress {
				msKeeper.SetIntermediaryAccountDelegator(ctx, imAddrressA, delA)
				msKeeper.SetIntermediaryAccountDelegator(ctx, imAddrressB, delB)
				return []sdk.AccAddress{delA, delB}
			},
			imAccs: []sdk.AccAddress{imAddrressA, imAddrressB},
			expPanic: false,
		},
		{
			name: "2 delegator, 2 intermediary account, failed",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.AccAddress {
				msKeeper.SetIntermediaryAccountDelegator(ctx, imAddrressA, delA)
				msKeeper.SetIntermediaryAccountDelegator(ctx, imAddrressA, delA)
				return []sdk.AccAddress{delA, delB}
			},
			imAccs: []sdk.AccAddress{imAddrressA, imAddrressB},
			expPanic: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
		
			if tc.expPanic {
				suite.Require().PanicsWithValue("intermediary account for delegator already set",func() {
					tc.malleate(suite.ctx, suite.msKeeper)
				})
			} else {
				inputs := tc.malleate(suite.ctx, suite.msKeeper)
				for idx, imAcc := range tc.imAccs {
					actualDel := suite.msKeeper.GetIntermediaryAccountDelegator(suite.ctx, imAcc)
					suite.Require().Equal(inputs[idx], actualDel)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSetDVPairSDKBondTokens() {
	delA := testutil.GenAddress()
	delB := testutil.GenAddress()
	valA := testutil.GenValAddress()
	valB := testutil.GenValAddress()

	bondSDKAmountA := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
	bondSDKAmountB := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200))

	testCases := []struct {
		name string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.Coin
		dels []sdk.AccAddress
		vals []sdk.ValAddress
		expPanic bool
	}{
		{
			name: "1 delegator, 1 validator, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.Coin {
				msKeeper.SetDVPairSDKBondTokens(ctx, delA, valA, bondSDKAmountA)
				return []sdk.Coin{bondSDKAmountA}
			},
			dels: []sdk.AccAddress{delA},
			vals: []sdk.ValAddress{valA},
			expPanic: false,
		},
		{
			name: "2 delegator, 2 validator, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.Coin {
				msKeeper.SetDVPairSDKBondTokens(ctx, delA, valA, bondSDKAmountA)
				msKeeper.SetDVPairSDKBondTokens(ctx, delB, valB, bondSDKAmountB)
				return []sdk.Coin{bondSDKAmountA, bondSDKAmountB}
			},
			dels: []sdk.AccAddress{delA, delB},
			vals: []sdk.ValAddress{valA, valB},
			expPanic: false,
		},
		{
			name: "1 delegator, 2 validator, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.Coin {
				msKeeper.SetDVPairSDKBondTokens(ctx, delA, valA, bondSDKAmountA)
				msKeeper.SetDVPairSDKBondTokens(ctx, delA, valB, bondSDKAmountB)
				return []sdk.Coin{bondSDKAmountA, bondSDKAmountB}
			},
			dels: []sdk.AccAddress{delA, delA},
			vals: []sdk.ValAddress{valA, valB},
			expPanic: false,
		},
		{
			name: "1 delegator, 1 validator, 2 bond amounts failed",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.Coin {
				msKeeper.SetDVPairSDKBondTokens(ctx, delA, valA, bondSDKAmountA)
				msKeeper.SetDVPairSDKBondTokens(ctx, delA, valA, bondSDKAmountB)
				return []sdk.Coin{bondSDKAmountB}
			},
			dels: []sdk.AccAddress{delA},
			vals: []sdk.ValAddress{valA},
			expPanic: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
		
			if tc.expPanic {
				suite.Require().PanicsWithValue("input token is not sdk bond token",func() {
					tc.malleate(suite.ctx, suite.msKeeper)
				})
			} else {
				inputs := tc.malleate(suite.ctx, suite.msKeeper)
				for idx, expOut := range inputs {
					actualCoin := suite.msKeeper.GetDVPairSDKBondTokens(suite.ctx, tc.dels[idx], tc.vals[idx])
					suite.Require().Equal(expOut, actualCoin)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSetDVPairBondTokens() {
	suite.SetupTest()

	delA := testutil.GenAddress()
	delB := testutil.GenAddress()
	valA := testutil.GenValAddress()
	valB := testutil.GenValAddress()

	bondAmountA := sdk.NewCoin("ario", sdk.NewInt(100))
	bondAmountB := sdk.NewCoin("arst", sdk.NewInt(200))

	suite.msKeeper.SetDVPairBondTokens(suite.ctx, delA, valA, bondAmountA)
	suite.msKeeper.SetDVPairBondTokens(suite.ctx, delB, valB, bondAmountB)

	suite.Equal(bondAmountA, suite.msKeeper.GetDVPairBondTokens(suite.ctx, delA, valA))
	suite.Equal(bondAmountB, suite.msKeeper.GetDVPairBondTokens(suite.ctx, delB, valB))

	suite.msKeeper.SetDVPairBondTokens(suite.ctx, delA, valB, bondAmountB)
	suite.Equal(bondAmountB, suite.msKeeper.GetDVPairBondTokens(suite.ctx, delA, valB))
}
