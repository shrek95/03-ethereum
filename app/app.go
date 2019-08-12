package main

import (
	"db/db/db"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type tx struct {
	TxID     int64
	BlockID  int64
	TxHash   string
	Value    int64
	Receiver string
	Sender   string
}

func main() {
	// multiplexer instance
	mux := httprouter.New()
	// mux /
	mux.GET("/", transaction)
	// mux /transactions
	mux.GET("/transactions", transactionDetails)

	// ListenAndServe for keeping server to listen on the port
	http.ListenAndServe(":8080", mux)
}

// transaction
func transaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "transaction.gohtml", nil)
	handleError(w, err)
}

// transactionDetails
func transactionDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("TransactionDetails func is executing :")

	// Parse the form
	err := r.ParseForm()
	// Handle error
	handleError(w, err)

	// create connection and get object for DB
	dtabse := db.CreateConn()
	// Choose DB
	db.ChooseDB(dtabse, "ethBlock")
	// Access the form data under the name address
	uaddr := r.FormValue("address")

	// get the slice of transaction for the specific address
	txSlice := db.GetTxFromUAddr(dtabse, uaddr)

	// Close db
	// db.CloseDB(dtabse)

	// execute the template populated with transaction data
	err = tpl.ExecuteTemplate(w, "transactiondetails.gohtml", txSlice)
	// Handle error
	handleError(w, err)
}

// HandleError ...
func handleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
