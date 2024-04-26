package staking

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	abci "github.com/cometbft/cometbft/abci/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModule embeds the Cosmos SDK's x/staking AppModuleBasic.
type AppModuleBasic struct {
	staking.AppModuleBasic
}

// AppModule embeds the Cosmos SDK's x/staking AppModule where we only override
// specific methods.
type AppModule struct {
	// embed the Cosmos SDK's x/staking AppModule
	staking.AppModule

	keeper     *keeper.Keeper
	accKeeper  types.AccountKeeper
	bankKeeper types.BankKeeper
}

// NewAppModule creates a new AppModule object using the native x/staking module
// AppModule constructor.
func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, subspace exported.Subspace) AppModule {
	stakingAppMod := staking.NewAppModule(cdc, keeper, ak, bk, subspace)
	return AppModule{
		AppModule:  stakingAppMod,
		keeper:     keeper,
		accKeeper:  ak,
		bankKeeper: bk,
	}
}

// InitGenesis delegates the InitGenesis call to the underlying x/staking module,
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)
	return am.keeper.InitGenesis(ctx, &genesisState)
}

// EndBlock delegates the EndBlock call to the underlying x/staking module,
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	updates := am.keeper.BlockValidatorUpdates(ctx)
	return updates
}
