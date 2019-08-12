package main

import (
	"context"
	"database/sql"
	"db/db/db"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var wg sync.WaitGroup

// transactions ...
func transactions(block *types.Block, client *ethclient.Client, dtabse *sql.DB) {

	i := 1 // count
	// Loop over transactions
	for _, tx := range block.Transactions() {
		fmt.Printf("------------Transaction count------------- %d \n", i)
		fmt.Printf("Hash of tx: ")
		fmt.Println(tx.Hash().Hex())
		fmt.Println("Value in the tx: ")
		fmt.Println(tx.Value().String())
		fmt.Println("To: ")
		fmt.Println(tx.To().Hex())
		// Get the chainID
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		// Get sender
		msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("From: ")
			fmt.Println(msg.From().Hex())
		}
		// ------------------------------------------------------

		// Convert value from big.Int to int
		valueBig := *(tx.Value())
		valueInt64 := valueBig.Int64()
		value := int(valueInt64)
		// convert block num to string
		blockNumBig := *(block.Number())
		blockNumInt := blockNumBig.Int64()
		blockNumStr := strconv.FormatInt(blockNumInt, 10)
		// insert block num, tx hash, value, to, from in table
		db.InsertTx(dtabse, blockNumStr, tx.Hash().Hex(), value, tx.To().Hex(), msg.From().Hex())
		// sleep for 3 sec
		duration := time.Duration(3) * time.Second
		time.Sleep(duration)
		i++
	}
	wg.Done()
}

// block ...
func block(dtabse *sql.DB) {

	// eth client
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	// get latest blockheader
	headerL, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// convert block number into int64 from *big.Int fromat
	lBlockBig := *(headerL.Number)
	lBlockInt64 := lBlockBig.Int64()
	// ------------------------------------------------------

	fmt.Println("Latest Block number :", headerL.Number.String())

	// get the starting block
	start := lBlockInt64 - 1

	// loop over the block range
	for i := start; i < lBlockInt64; i++ {
		fmt.Printf("------------------Block Number-------------------- %d \n", i)
		// create object of type big.Int
		holder := new(big.Int)
		// convert the value in i to *big.Int and set it in holder
		holder.SetString(strconv.Itoa(int(i)), 10)

		// Get block number in var block
		block, err := client.BlockByNumber(context.Background(), holder)
		if err != nil {
			log.Fatal(err)
		}
		// Get transaction count
		count, err := client.TransactionCount(context.Background(), block.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Transaction count: %d ", count)
		// Convert block num into string
		blockNumBig := *(block.Number())
		blockNumInt := blockNumBig.Int64()
		blockNumStr := strconv.FormatInt(blockNumInt, 10)

		// Insert blocknumber, block hash, tx count in the block
		db.InsertBlock(dtabse, blockNumStr, block.Hash().Hex(), int(count))

		// sleep for 2 sec
		duration := time.Duration(2) * time.Second
		time.Sleep(duration)
		// ---------------------------------------

		// initiate transaction func on a go routine
		wg.Add(1)
		go transactions(block, client, dtabse)
	}
	wg.Wait()
	// ---------------------------------------------
}

func main() {
	fmt.Println("Block entry start: ")
	// create connection with the database
	dtabse := db.CreateConn()
	// choose the database
	db.ChooseDB(dtabse, "ethBlock")
	// db.DeleteTable(dtabse, "ethBlock", "transaction")
	// db.DeleteTable(dtabse, "ethBlock", "block")
	// db.CreateTable(dtabse)

	block(dtabse)
	db.CloseDB(dtabse)

}
