package main

import (
	"db/db/db"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type tx struct {
	Tx_id    int64
	Block_id int64
	Tx_hash  string
	Value    int64
	Receiver string
	Sender   string
}

func main() {
	mux := httprouter.New()
	mux.GET("/", transaction)
	mux.POST("/transactions/", transactionDetails)
	http.ListenAndServe(":8080", mux)
}

func transaction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Println("transaction func triggered:")
	err := tpl.ExecuteTemplate(w, "transaction.gohtml", nil)
	HandleError(w, err)
}

func transactionDetails(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)
	fmt.Println("transactionDetails func triggered:")
	dtabse := db.CreateConn()
	db.ChooseDB(dtabse, "ethBlock")
	uaddr := req.Form["address"]
	fmt.Println("uaddr type : ", reflect.TypeOf(uaddr))
	addr := uaddr[0]
	tx_slice := db.GetTxFromUAddr(dtabse, addr /*"0x34371D7Bd50936CB478145f11F7a3B24bf9D9D92"*/)
	fmt.Println("tx_slice type", reflect.TypeOf(tx_slice))
	err := tpl.ExecuteTemplate(w, "transactiondetails.gohtml", tx_slice)
	HandleError(w, err)
}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
