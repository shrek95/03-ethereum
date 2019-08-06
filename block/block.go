package main

import (
	"context"
	"database/sql"
	"db/db/db"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var wg sync.WaitGroup

func transactions(block *types.Block, client *ethclient.Client, dtabse *sql.DB) {

	i := 1
	for _, tx := range block.Transactions() {
		fmt.Printf("------------Transaction count------------- %d \n", i)
		fmt.Println(tx.Hash().Hex())
		fmt.Println(tx.Value().String())
		fmt.Println(tx.To().Hex())

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(msg.From().Hex())
		}

		value_big := *(tx.Value())
		value_int64 := value_big.Int64()
		value := int(value_int64)

		block_num_big := *(block.Number())
		block_num_int := block_num_big.Int64()
		block_num_str := strconv.FormatInt(block_num_int, 10)

		db.InsertTx(dtabse, block_num_str, tx.Hash().Hex(), value, tx.To().Hex(), msg.From().Hex())
		i++
	}
	wg.Done()
}

func block(dtabse *sql.DB) {

	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	headerL, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	lBlockBig := *(headerL.Number)
	lBlockInt64 := lBlockBig.Int64()
	fmt.Println(reflect.TypeOf(lBlockInt64))
	fmt.Println(reflect.TypeOf(headerL.Number))
	start := lBlockInt64 - 10000

	for i := start; i < lBlockInt64; i++ {
		fmt.Printf("-------Block Number----- %d \n", i)
		holder := new(big.Int)
		holder.SetString(strconv.Itoa(int(i)), 10)

		block, err := client.BlockByNumber(context.Background(), holder)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Block number")
		fmt.Println(block.Number) //func() big.Int
		fmt.Println("Block hash")
		fmt.Println(block.Hash().Hex()) // string
		fmt.Println("transaction array len")
		fmt.Println(len(block.Transactions()))
		fmt.Println("Transaction count using function")
		count, err := client.TransactionCount(context.Background(), block.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(count)
		block_num_big := *(block.Number())
		block_num_int := block_num_big.Int64()
		block_num_str := strconv.FormatInt(block_num_int, 10)
		db.InsertBlock(dtabse, block_num_str, block.Hash().Hex(), int(count))
		wg.Add(1)
		go transactions(block, client, dtabse)
	}
	wg.Wait()
}

func main() {
	fmt.Println("Here we go, main is triggered: ")
	dtabse := db.CreateConn()
	db.ChooseDB(dtabse, "ethBlock")
	//	db.DeleteTable(dtabse, "ethBlock", "transaction")
	//	db.DeleteTable(dtabse, "ethBlock", "block")
	//	db.CreateTable(dtabse)
	block(dtabse)

}
