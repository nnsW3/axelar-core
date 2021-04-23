package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	ethereumTypes "github.com/axelarnetwork/axelar-core/x/ethereum/types"
)

const (
	flagGateway  = "gateway"
	flagToken    = "token"
	flagBurnable = "burnable"
)

// SetGenesisEthContractsCmd returns set-genesis-chain-params cobra Command.
func SetGenesisEthContractsCmd(defaultNodeHome string) *cobra.Command {
	var gatewayFile, tokenFile, burnableFile string

	cmd := &cobra.Command{
		Use:   "set-genesis-ethereum-contracts",
		Short: "Set the ethereum's contract parameters in genesis.json",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			genesisState := ethereumTypes.GetGenesisStateFromAppState(cdc, appState)

			if gatewayFile != "" {
				gateway, err := getByteCodes(gatewayFile)
				if err != nil {
					return err
				}
				genesisState.Params.Gateway = gateway
			}

			if tokenFile != "" {
				token, err := getByteCodes(tokenFile)
				if err != nil {
					return err
				}
				genesisState.Params.Token = token
			}

			if burnableFile != "" {
				burnable, err := getByteCodes(burnableFile)
				if err != nil {
					return err
				}
				genesisState.Params.Burnable = burnable
			}

			genesisStateBz, err := cdc.MarshalJSON(&genesisState)
			if err != nil {
				return fmt.Errorf("failed to marshal ethereum genesis state: %w", err)
			}
			appState[ethereumTypes.ModuleName] = genesisStateBz
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")

	cmd.Flags().StringVar(&gatewayFile, flagGateway, "", "Path to the Axelar Gateway contract ABI.")
	cmd.Flags().StringVar(&tokenFile, flagToken, "", "Path to the tokens contract ABI.")
	cmd.Flags().StringVar(&burnableFile, flagBurnable, "", "Path to the burner contract ABI.")

	return cmd
}