package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type transaction struct {
	tx_id    int64
	block_id int64
	tx_hash  string
	value    int
	receiver string
	sender   string
}

type block struct {
	block_id   int64
	block_num  string
	block_hash string
	tx_count   int
}

var ip = os.GetEnv("DOCKER_IP")

var path = "root:root@tcp(" + ip + ":3306)/"

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

func CreateDB(db *sql.DB, name string) {
	query := "CREATE DATABASE " + name
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Database created successfully: ")
	}
}

func ChooseDB(db *sql.DB, dbName string) {
	query := "USE " + dbName
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("DB selected successfully: " + dbName)
	}
}

func CreateTable(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS block( block_id bigint NOT NULL AUTO_INCREMENT, block_num varchar(200), block_hash varchar(200), tx_count int, PRIMARY KEY (block_id) );")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Block Table created successfully..")
	}
	stmttx, err := db.Prepare("CREATE TABLE IF NOT EXISTS transaction(tx_id bigint NOT NULL AUTO_INCREMENT, block_id bigint NOT NULL, tx_hash varchar(200), value bigint, receiver varchar(100), sender varchar(100), PRIMARY KEY (tx_id), FOREIGN KEY (block_id) REFERENCES block(block_id) );")
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

// value ko bigint karna table creation me
func InsertTx(db *sql.DB, block_num string, tx_hash string, value int, receiver string, sender string) {
	blockQuery := " SELECT block.block_id FROM block WHERE block.block_num = ?;"
	rows, err := db.Query(blockQuery, block_num)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	var block_id int
	for rows.Next() {
		if err := rows.Scan(&block_id); err != nil {
			fmt.Println(err.Error())
		}
	}

	stmt, err := db.Prepare("INSERT INTO transaction( block_id, tx_hash, value, receiver, sender) VALUES( ?, ?, ?, ?, ? )")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec(block_id, tx_hash, value, receiver, sender)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Entry is successfull in transaction table: ")
	}
}

func InsertBlock(db *sql.DB, block_num string, block_hash string, tx_count int) {
	stmt, err := db.Prepare("INSERT INTO block(block_num, block_hash, tx_count) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Preparation successfull for block insert: ")
	}

	_, err = stmt.Exec(block_num, block_hash, tx_count)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Entry is block table is successfull: ")
	}
}
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

func GetTxOfBlockNumQuery(db *sql.DB, tName string, block_num int64) []transaction {
	query := "SELECT transaction.* FROM block JOIN transaction ON block.block_id = transaction.block_id WHERE block.block_num = " + strconv.Itoa(int(block_num)) + ";"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Here is the error: " + err.Error())
	} else {
		fmt.Println("Query Preparation Successfull:")
	}

	defer rows.Close()

	txSlice := make([]transaction, 0)

	for rows.Next() {
		var tx transaction
		if err := rows.Scan(&tx.tx_id, &tx.block_id, &tx.tx_hash, &tx.value, &tx.receiver, &tx.sender); err != nil {
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

func GetTxFromUAddr(db *sql.DB, uAddr string) []transaction {
	fmt.Println("HEre is the variable printed ______", uAddr)
	query := "SELECT * FROM transaction WHERE transaction.sender = " + uAddr + " OR " + "transaction.receiver = " + uAddr + ";"
	// query := "SELECT * FROM transaction;"
	// select * from transaction where transaction.sender = '0x662710a199415B48B210F8dc9937083526e12583' or transaction.receiver = '0x34371D7Bd50936CB478145f11F7a3B24bf9D9D92'
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Here is the error: " + err.Error())
	} else {
		fmt.Println("Query Preparation Successfull:")
	}

	defer rows.Close()

	txSlice := make([]transaction, 0)

	for rows.Next() {
		var tx transaction
		if err := rows.Scan(&tx.tx_id, &tx.block_id, &tx.tx_hash, &tx.value, &tx.receiver, &tx.sender); err != nil {
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

func CloseDB(db *sql.DB) {
	db.Close()
}

/*
func main() {
	db := CreateConn()

	// CreateDB(db, "ethBlock")
	ChooseDB(db, "ethBlock")
	CreateTable(db)
	//	DeleteTable(db, "ethBlock", "transaction")
	//	DeleteTable(db, "ethBlock", "block")
	//	Insert(db, 223344, "0x223344", "0x556677", 12)
	//	Insert(db, 223345, "0x223355", "0x556688", 13)
	// 	InsertTx(db, 112233, 223344, "0x998899", 100000, "0x787979", "0x3953000")

	// InsertTx(db, 112234, 223344, "0x998898", 10, "0x7879", "0x395")
	// InsertTx(db, 112235, 223344, "0x998897", 10, "0x7877", "0x3956")
	//	GetTxOfBlockIdQuery(db, "block", 223344)
	defer CloseDB(db)

}*/
