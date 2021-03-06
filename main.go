package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"database/sql"
	"github.com/DavidHuie/gomigrate"
	"fmt"
)

func main() {
	ApplyMigrations()

	router := mux.NewRouter()
	router.HandleFunc("/app/expenses", CreateExpenseEndpoint).Methods("POST")
	router.HandleFunc("/app/expenses", GetExpenseEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}


func CreateExpenseEndpoint(w http.ResponseWriter, req *http.Request) {
	var expense Expense
	_ = json.NewDecoder(req.Body).Decode(&expense)

	SaveExpense(expense)
	json.NewEncoder(w).Encode(expense)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", "POST, GET")
	w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
	w.Header().Set("Access-Control-Max-Age", 86400)

}

func GetExpenseEndpoint(w http.ResponseWriter, req *http.Request) {
	expenses := GetAllExpenses()
	_ = json.NewDecoder(req.Body).Decode(&expenses)
	json.NewEncoder(w).Encode(expenses)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", "POST, GET")
	w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
	w.Header().Set("Access-Control-Max-Age", 86400)
}


type Expense struct {
	ID     int32  `json:"id,omitempty"`
	Date   string `json:"date,omitempty"`
	Reason string `json:"reason,omitempty"`
	Amount string `json:"amount"`
	VAT    string `json:"vat"`
}

func SaveExpense(expense Expense) {
	db, _ := sql.Open("mysql", "userDev:passDev@tcp(127.0.0.1:3306)/expenses2?parseTime=true")
	defer db.Close()
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
}

func GetAllExpenses() []Expense{
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
	return expenses
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
