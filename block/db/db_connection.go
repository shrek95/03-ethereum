package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

// Transaction ...
type Transaction struct {
	txID     int64
	blockID  int64
	txHash   string
	value    int
	receiver string
	sender   string
}

// Block ...
type Block struct {
	blockID   int64
	blockNum  string
	blockHash string
	txCount   int
}

var ip = os.Getenv("DOCKER_IP")
var db = os.Getenv("DB_NAME")

var path = "root:root@tcp(" + ip + ":3306)/" + db

// CreateConn ...
func CreateConn() *sql.DB {
	db, err := sql.Open("mysql", path)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Connection created successfully: ")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Connection not working: ")
	} else {
		fmt.Println("Connection working perfectly: ")
	}

	return db
}

// CreateDB ...
func CreateDB(db *sql.DB, name string) {
	query := "CREATE DATABASE " + name
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Database created successfully: ")
	}
}

// ChooseDB ...
func ChooseDB(db *sql.DB, dbName string) {
	query := "USE " + dbName
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("DB selected successfully: " + dbName)
	}
}

// CreateTable ...
func CreateTable(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS block( blockID bigint NOT NULL AUTO_INCREMENT, blockNum varchar(200), blockHash varchar(200), txCount int, PRIMARY KEY (blockID) );")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Block Table created successfully..")
	}
	stmttx, err := db.Prepare("CREATE TABLE IF NOT EXISTS transaction(txID bigint NOT NULL AUTO_INCREMENT, blockID bigint NOT NULL, txHash varchar(200), value bigint, receiver varchar(100), sender varchar(100), PRIMARY KEY (txID), FOREIGN KEY (blockID) REFERENCES block(blockID) );")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmttx.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Transaction Table created successfully..")
	}
}

// InsertTx ...
func InsertTx(db *sql.DB, blockNum string, txHash string, value int, receiver string, sender string) {
	blockQuery := " SELECT block.blockID FROM block WHERE block.blockNum = ?;"
	rows, err := db.Query(blockQuery, blockNum)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	var blockID int
	for rows.Next() {
		if err := rows.Scan(&blockID); err != nil {
			fmt.Println(err.Error())
		}
	}

	stmt, err := db.Prepare("INSERT INTO transaction( blockID, txHash, value, receiver, sender) VALUES( ?, ?, ?, ?, ? )")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec(blockID, txHash, value, receiver, sender)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Entry in successfull in transaction table: ")
	}
}

// InsertBlock ...
func InsertBlock(db *sql.DB, blockNum string, blockHash string, txCount int) {
	stmt, err := db.Prepare("INSERT INTO block(blockNum, blockHash, txCount) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Preparation successfull for block insert: ")
	}

	_, err = stmt.Exec(blockNum, blockHash, txCount)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Entry is block table is successfull: ")
	}
}

// DeleteTable ...
func DeleteTable(db *sql.DB, dbName string, tName string) {

	dropTable := "DROP TABLE " + tName
	stmt, err := db.Prepare(dropTable)
	if err != nil {
		fmt.Println("Here is the error:      " + err.Error())
	} else {
		fmt.Println("Delete Preparation Successfull for : " + tName)
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("Here is the error: " + err.Error())
	} else {
		fmt.Println("Delete Successful of " + tName)
	}
}

// GetTxOfBlockNumQuery ...
func GetTxOfBlockNumQuery(db *sql.DB, tName string, blockNum int64) []Transaction {
	query := "SELECT transaction.* FROM block JOIN transaction ON block.blockID = transaction.blockID WHERE block.blockNum = " + strconv.Itoa(int(blockNum)) + ";"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Here is the error: " + err.Error())
	} else {
		fmt.Println("Query Preparation Successfull:")
	}

	defer rows.Close()

	txSlice := make([]Transaction, 0)

	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.txID, &tx.blockID, &tx.txHash, &tx.value, &tx.receiver, &tx.sender); err != nil {
			fmt.Println("Row scan failed: ")
		}
		txSlice = append(txSlice, tx)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Query is successful \nHere is the data:  ")
		fmt.Println(txSlice)
	}
	return txSlice
}

// GetTxFromUAddr ...
func GetTxFromUAddr(db *sql.DB, uAddr string) []Transaction {
	query := "SELECT * FROM transaction WHERE transaction.sender=" + uAddr + " OR " + "transaction.receiver=" + uAddr + ";"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Here is the error: " + err.Error())
	} else {
		fmt.Println("Query Preparation Successfull:")
	}

	defer rows.Close()

	txSlice := make([]Transaction, 0)

	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.txID, &tx.blockID, &tx.txHash, &tx.value, &tx.receiver, &tx.sender); err != nil {
			fmt.Println("Row scan failed: ")
		}
		txSlice = append(txSlice, tx)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Query is successful \nHere is the data:  ")
		fmt.Println(txSlice)
	}
	return txSlice
}

// CloseDB ...
func CloseDB(db *sql.DB) {
	db.Close()
}
