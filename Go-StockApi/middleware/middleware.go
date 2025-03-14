package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stock/config"
	"stock/models"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message.omitempty"`
}

func connectDatabase() *sql.DB {

	DB_HOST := config.GetEnv("DB_HOST")
	DB_PORT := config.GetEnv("DB_PORT")
	DB_USER := config.GetEnv("DB_USER")
	DB_PASSWORD := config.GetEnv("DB_PASSWORD")
	DB_NAME := config.GetEnv("DB_NAME")
	dsn_info := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dsn_info)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected")

	return db
}
func CreateStock(w http.ResponseWriter, r *http.Request) {

	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request Body. %v", err)
	}

	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "Stock Created Successfully",
	}

	json.NewEncoder(w).Encode(res)

}

func GetStock(w http.ResponseWriter, r *http.Request) {
	parmas := mux.Vars(r)

	id, err := strconv.Atoi(parmas["id"])
	if err != nil {
		log.Fatalf("unable to convert String to INT, %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get strock %v", err)
	}
	json.NewEncoder(w).Encode((stock))

}

func GetAllStock(w http.ResponseWriter, r *http.Request) {

	stocks, err := getAllstocks()

	if err != nil {
		log.Fatalf("Unable to fetch all stocks. %v ", err)
	}
	json.NewEncoder(w).Encode(stocks)

}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to into int. %v", err)
	}
	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode Request Body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock Updated successfully. Total row/ record affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string to integer %v", err)
	}
	deleterRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Successfully updated stock . Total Rows/record Affected %v", deleterRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

//-----------------------HANDLER FUNCTIONS------------------

func insertStock(stock models.Stock) int64 {

	db := connectDatabase()
	defer db.Close()

	sqlStatement := "INSERT INTO STOCK  (name,price,company) VALUES ($1 ,$2, $3)RETURNING stockid"
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Unable to insert stock %v", err)
	}
	fmt.Printf("Inserted Single Record. %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := connectDatabase()
	defer db.Close()

	var stock models.Stock
	sqlStament := " SELECT * FROM stock WHERE stockid = $1"
	row := db.QueryRow(sqlStament, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to see scan row. %v", err)
	}
	return stock, err

}

func getAllstocks() ([]models.Stock, error) {

	db := connectDatabase()
	defer db.Close()

	var stocks []models.Stock

	sqlStatemet := "SELECT * FROM stocks"

	rows, err := db.Query(sqlStatemet)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("unableto scan row. %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err

}
func updateStock(id int64, stock models.Stock) int64 {
	db := connectDatabase()
	defer db.Close()
	sqlStatement := "UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1"

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execure query %v", err)
	}
	rowAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while chekcing affected rows %v", err)
	}
	fmt.Printf("Total rows affected %v ", err)

	return rowAffected
}

func deleteStock(id int64) int64 {
	db := connectDatabase()
	defer db.Close()

	sqlStatement := "DELETE FROM stocks WHERE stockid=$1"

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Error during executing Query %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected

}
