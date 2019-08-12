package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

// Transaction ... transaction structure
type Transaction struct {
	TxID     int64
	BlockID  int64
	TxHash   string
	Value    int
	Receiver string
	Sender   string
}

// Block ... block structure
type Block struct {
	BlockID   int64
	BlockNum  string
	BlockHash string
	TxCount   int
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
		fmt.Println("Connection not working: " + err.Error())
	} else {
		fmt.Println("Connection working perfectly: ")
	}

	return db
}

// CreateDB ..,
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
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS block( BlockID bigint NOT NULL AUTO_INCREMENT, BlockNum varchar(200), BlockHash varchar(200), TxCount int, PRIMARY KEY (BlockID) );")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Block Table created successfully..")
	}
	stmttx, err := db.Prepare("CREATE TABLE IF NOT EXISTS transaction(TxID bigint NOT NULL AUTO_INCREMENT, BlockID bigint NOT NULL, TxHash varchar(200), Value bigint, Receiver varchar(100), Sender varchar(100), PRIMARY KEY (TxID), FOREIGN KEY (BlockID) REFERENCES block(BlockID) );")
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
func InsertTx(db *sql.DB, BlockNum string, TxHash string, Value int, Receiver string, Sender string) {
	blockQuery := " SELECT block.BlockID FROM block WHERE block.BlockNum = ?;"
	rows, err := db.Query(blockQuery, BlockNum)
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

	stmt, err := db.Prepare("INSERT INTO transaction( BlockID, TxHash, Value, Receiver, Sender) VALUES( ?, ?, ?, ?, ? )")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec(blockID, TxHash, Value, Receiver, Sender)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Entry is successfull in transaction table: ")
	}
}

// InsertBlock ...
func InsertBlock(db *sql.DB, BlockNum string, BlockHash string, TxCount int) {
	stmt, err := db.Prepare("INSERT INTO block(BlockNum, BlockHash, TxCount) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Preparation successfull for block insert: ")
	}

	_, err = stmt.Exec(BlockNum, BlockHash, TxCount)
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
func GetTxOfBlockNumQuery(db *sql.DB, tName string, BlockNum int64) []Transaction {
	query := "SELECT transaction.* FROM block JOIN transaction ON block.BlockID = transaction.BlockID WHERE block.BlockNum = ? ;"
	rows, err := db.Query(query, strconv.Itoa(int(BlockNum)))
	if err != nil {
		fmt.Println("Here is the error: " + err.Error())
	} else {
		fmt.Println("Query Preparation Successfull:")
	}

	defer rows.Close()

	txSlice := make([]Transaction, 0)

	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.TxID, &tx.BlockID, &tx.TxHash, &tx.Value, &tx.Receiver, &tx.Sender); err != nil {
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

// GetTxFromUAddr ... gets the transaction linked with the User address
func GetTxFromUAddr(db *sql.DB, uAddr string) []Transaction {
	fmt.Println("HEre is the variable printed ______" + uAddr)
	// query := "SELECT * FROM transaction WHERE transaction.sender = 0x32A4FdD43eDd3319F7941b69924A24AC0EF501c2" + ";" /*+ " AND " + "transaction.receiver=" + uAddr + ");"*/
	// query := "SELECT * FROM transaction;"
	// select * from transaction where transaction.sender = '0x662710a199415B48B210F8dc9937083526e12583' or transaction.receiver = '0x34371D7Bd50936CB478145f11F7a3B24bf9D9D92'
	query := `SELECT * FROM transaction WHERE transaction.sender = ? OR transaction.receiver = ?;`
	// rows, err := db.Query(query, "0x32A4FdD43eDd3319F7941b69924A24AC0EF501c2", "0x32A4FdD43eDd3319F7941b69924A24AC0EF501c2")
	rows, err := db.Query(query, uAddr, uAddr)
	if err != nil {
		fmt.Println("Here is the error: " + err.Error())
	} else {
		fmt.Println("Query Preparation Successfull:")
	}

	defer rows.Close()

	txSlice := make([]Transaction, 0)

	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.TxID, &tx.BlockID, &tx.TxHash, &tx.Value, &tx.Receiver, &tx.Sender); err != nil {
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

// CloseDB ... closes the db
func CloseDB(db *sql.DB) {
	db.Close()
}

/*
func main() {
	db := CreateConn()

	// CreateDB(db, "ethBlock")
	ChooseDB(db, "ethBlock")
	CreateTable(db)
	// DeleteTable(db, "ethBlock", "transaction")
	// DeleteTable(db, "ethBlock", "block")
	//	Insert(db, 223344, "0x223344", "0x556677", 12)
	//	Insert(db, 223345, "0x223355", "0x556688", 13)
	// 	InsertTx(db, 112233, 223344, "0x998899", 100000, "0x787979", "0x3953000")

	// InsertTx(db, 112234, 223344, "0x998898", 10, "0x7879", "0x395")
	// InsertTx(db, 112235, 223344, "0x998897", 10, "0x7877", "0x3956")
	//	GetTxOfBlockIdQuery(db, "block", 223344)
	defer CloseDB(db)

}*/
