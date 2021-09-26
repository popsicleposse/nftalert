package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rbrick/nftalert/contract/erc721meta"
	ethutil "github.com/rbrick/nftalert/utils"
)

/*

What is needed?


- A way to read the smart contract & tell what it is. This will be done within the contract package.
  I figure, we can probably use a few free etherscan API keys for this. Even having just 3 is 300,000 queries
  per day, which should be more than plenty to collect information. What we'd do is when we come across a transaction
  from a tracked account, we can use etherscan to either read the contract, or simply grab the ABI, and we then store
  the ABI & dump that in a database somewhere. I'd like not to rely too heavily on something centralized. But it's by far the
  easiest way.

  Another thought was using some Solidity bytecode decompiler, or just scanning the bytecode disassembly for signatures. We could implement a
  few different strategies.

  but the basic idea, since a contract would mostly be unchanging, especially for our ERC721 target, we can grab info on it,
  and cache it indefinitely.

  We'll probably need to use a mix of strategies if we're really wanting to get the most out of

- A database. Probably MongoDB, or some SQL server like Postgres. Should be pretty simple to model. Think more about the database modeling.
 I'd like for the database to store every interaction (and previous interactions), so we can get graphs and charts going

- The bot/site itself. Can be some React/Vue.js site that just basically lists incoming transactions from our watchlist. Possibly a page
   for adding/removing address watchlists. This shouldn't take too much time.


Mostly the trick is to learn & understand more about Ethereum's ABI, it's all in all pretty simple to understand & work with,
seems like it'd be pretty smooth sailing from there. I think i'll work with and explore ether's ABI related commands in geth

found here

https://github.com/ethereum/go-ethereum/tree/master/cmd/abidump


and here https://github.com/ethereum/go-ethereum/tree/master/cmd/abigen

There is also an entire EVM command program which seems like it might be a treasure trove of useful examples

https://github.com/ethereum/go-ethereum/tree/master/cmd/evm


*/

func main() {

	// the websocket URI can be anything, something like ether.rcw.io
	// which can then actually be a load balancer to like 10 separate RPC nodes.
	// Ideally you'd have full nodes. Or maybe like one or two full nodes that store a
	// full copy. Then light nodes that connect solely to those full nodes as peers
	// and then load balance between the light nodes which are relaying between the full nodes
	// I think that'll provide affordability & stability.
	rpcClient, err := rpc.Dial("ws://localhost:8546")
	if err != nil {
		log.Fatalln(err)
	}
	client := ethclient.NewClient(rpcClient)
	// logsChannel := make(chan types.Log)
	// client.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{}, logsChannel)

	// for l := range logsChannel {
	// 	fmt.Println(l)
	// }

	headChannel := make(chan *types.Header)

	client.SubscribeNewHead(context.Background(), headChannel)

	for {
		head := <-headChannel
		block, err := client.BlockByHash(context.Background(), head.ParentHash)

		if err != nil {
			log.Fatalln(err)
		}

		for idx, transaction := range block.Body().Transactions {

			receipt, err := client.TransactionReceipt(context.Background(), transaction.Hash())

			if receipt == nil {

				fmt.Println(err)
			}

			from, _ := ethutil.GetSenderAddress(rpcClient, block.Number(), uint(idx))
			// js, _ := transaction.MarshalJSON()

			if transaction.To() != nil {
				smartContract := false
				code, err := client.CodeAt(context.Background(), *transaction.To(), block.Number())

				if err == nil && len(code) > 0 {
					smartContract = true

					erc721meta, err := erc721meta.NewErc721meta(*transaction.To(), client)

					if err == nil {
						name, err := erc721meta.Name(&bind.CallOpts{})
						if err == nil {
							fmt.Println("Token Name:", name)
						}
					}
				}

				fmt.Println(from.String(), "->", transaction.To(), ethutil.FormatValue(transaction.Value()), "eth", "is smart contract?", smartContract)
			} else {
				fmt.Println(from.String(), "-> 0x0 (Unknown)", transaction.Value(), "eth")
			}
		}
		// break
	}

}
