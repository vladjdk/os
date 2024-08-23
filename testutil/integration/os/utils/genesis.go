// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package utils

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	exampleapp "github.com/evmos/os/example_chain"
	"github.com/evmos/os/testutil"
	testkeyring "github.com/evmos/os/testutil/integration/os/keyring"
	"github.com/evmos/os/testutil/integration/os/network"
	erc20types "github.com/evmos/os/x/erc20/types"
)

const (
	// erc20TokenPairHex is the string representation of the ERC-20 token pair address.
	erc20TokenPairHex = "0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd" //#nosec G101 -- these are not hardcoded credentials #gitleaks:allow
)

func CreateGenesisWithTokenPairs(keyring testkeyring.Keyring) network.CustomGenesisState {
	// Add all keys from the keyring to the genesis accounts as well.
	//
	// NOTE: This is necessary to enable the account to send EVM transactions,
	// because the Mono ante handler checks the account balance by querying the
	// account from the account keeper first. If these accounts are not in the genesis
	// state, the ante handler finds a zero balance because of the missing account.
	accs := keyring.GetAllAccAddrs()
	genesisAccounts := make([]*authtypes.BaseAccount, len(accs))
	for i, addr := range accs {
		genesisAccounts[i] = &authtypes.BaseAccount{
			Address:       addr.String(),
			PubKey:        nil,
			AccountNumber: uint64(i + 1),
			Sequence:      1,
		}
	}

	accGenesisState := authtypes.DefaultGenesisState()
	for _, genesisAccount := range genesisAccounts {
		// NOTE: This type requires to be packed into a *types.Any as seen on SDK tests,
		// e.g. https://github.com/evmos/cosmos-sdk/blob/v0.47.5-evmos.2/x/auth/keeper/keeper_test.go#L193-L223
		accGenesisState.Accounts = append(accGenesisState.Accounts, codectypes.UnsafePackAny(genesisAccount))
	}

	// Add token pairs to genesis
	erc20GenesisState := exampleapp.NewErc20GenesisState()
	erc20GenesisState.TokenPairs = append(erc20GenesisState.TokenPairs,
		erc20types.TokenPair{
			Erc20Address:  erc20TokenPairHex,
			Denom:         "xmpl",
			Enabled:       true,
			ContractOwner: erc20types.OWNER_MODULE, // NOTE: Owner is the module account since it's a native token and was registered through governance
		},
		erc20types.TokenPair{
			Erc20Address:  testutil.WEVMOSContractTestnet,
			Denom:         testutil.ExampleAttoDenom,
			Enabled:       true,
			ContractOwner: erc20types.OWNER_MODULE, // NOTE: Owner is the module account since it's a native token and was registered through governance
		},
	)

	// Combine module genesis states
	return network.CustomGenesisState{
		authtypes.ModuleName:  accGenesisState,
		erc20types.ModuleName: erc20GenesisState,
	}
}
