package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

func GetSenderAddress(client *rpc.Client, block *big.Int, idx uint) (*common.Address, error) {
	var meta struct {
		Hash common.Hash
		From common.Address
	}

	err := client.Call(&meta, "eth_getTransactionByBlockNumberAndIndex", hexutil.EncodeBig(block), hexutil.Uint64(idx).String())

	if err != nil {
		return nil, err
	}
	return &meta.From, nil
}
