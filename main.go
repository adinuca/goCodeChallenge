package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	ApplyMigrations()

	router := mux.NewRouter()
	router.HandleFunc("/app/expenses", CreateExpenseEndpoint).Methods("POST")
	router.HandleFunc("/app/expenses", GetExpenseEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}

func ApplyMigrations() {
	db, _ := sql.Open("mysql", "userDev:passDev@tcp(127.0.0.1:3306)/expenses2?parseTime=true")
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	migrator, _ := gomigrate.NewMigrator(db, gomigrate.Mysql{}, "/Users/adina/go/src/goCodeChallenge/migrations")
	err := migrator.Migrate()
	if err != nil {
		fmt.Println(err.Error())
	}
}

type Expense struct {
	ID     int32  `json:"id,omitempty"`
	Date   string `json:"date,omitempty"`
	Reason string `json:"reason,omitempty"`
	Amount string `json:"amount"`
	VAT    string `json:"vat"`
}

func CreateExpenseEndpoint(w http.ResponseWriter, req *http.Request) {
	db, _ := sql.Open("mysql", "userDev:passDev@tcp(127.0.0.1:3306)/expenses2?parseTime=true")
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	var expense Expense
	_ = json.NewDecoder(req.Body).Decode(&expense)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	fmt.Println(expense)
	stmt, err := tx.Prepare("INSERT INTO expense (date,amount, reason) VALUES (?, ?, ?)")
	amount, err := strconv.ParseFloat(expense.Amount, 64)
	_, err = stmt.Exec(expense.Date, amount, expense.Reason)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(expense)
}

func GetExpenseEndpoint(w http.ResponseWriter, req *http.Request) {
	db, _ := sql.Open("mysql", "userDev:passDev@tcp(127.0.0.1:3306)/expenses2?parseTime=true")
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT date, amount, reason FROM expense")
	if err != nil {
		log.Fatal(err)
	}
	var expenses []Expense = []Expense{}
	for rows.Next() {
		var date, reason string
		var amount float64
		if err := rows.Scan(&date, &amount, &reason); err != nil {
			log.Fatal(err)
		}
		amountString := strconv.FormatFloat(amount, 'f', 6, 64)
		vatString := strconv.FormatFloat(amount*20/100, 'f', 6, 64)
		expense := Expense{Date: date, Reason: reason, Amount: amountString, VAT: vatString}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	//var expense Expense = Expense{ID:1, Date: "12/12/1988", Reason:"just so", Amount:12, VAT: 12 * 20 / 100 }
	_ = json.NewDecoder(req.Body).Decode(&expenses)
	json.NewEncoder(w).Encode(expenses)
}
