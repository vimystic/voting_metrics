package cosmos

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func defaultPageRequest() *querytypes.PageRequest {
	return &querytypes.PageRequest{
		Key:        []byte(""),
		Offset:     0,
		Limit:      1000,
		CountTotal: false,
	}
}

// QueryABCI performs an ABCI query and returns the appropriate response and error sdk error code.
func (cc *CosmosProvider) QueryABCI(ctx context.Context, req abci.RequestQuery) (abci.ResponseQuery, error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: req.Height,
		Prove:  req.Prove,
	}
	result, err := cc.RPCClient.ABCIQueryWithOptions(ctx, req.Path, req.Data, opts)
	if err != nil {
		return abci.ResponseQuery{}, err
	}

	if !result.Response.IsOK() {
		return abci.ResponseQuery{}, sdkErrorToGRPCError(result.Response)
	}

	return result.Response, nil
}

// queryBalance returns the amount of coins in the relayer account with address as input
func (cc *CosmosProvider) QueryBalance(ctx context.Context, address string, height uint64) (sdk.Coins, error) {
	qc := banktypes.NewQueryClient(cc)
	p := defaultPageRequest()
	coins := sdk.Coins{}

	headers := make(map[string]string)

	if height > 0 {
		headers[grpctypes.GRPCBlockHeightHeader] = fmt.Sprintf("%d", height)
	}
	md := metadata.New(headers)

	for {
		res, err := qc.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
			Address:    address,
			Pagination: p,
		}, grpc.Header(&md))
		if err != nil {
			return nil, err
		}

		coins = append(coins, res.Balances...)
		next := res.GetPagination().GetNextKey()
		if len(next) == 0 {
			break
		}

		p.Key = next
	}
	return coins, nil
}
