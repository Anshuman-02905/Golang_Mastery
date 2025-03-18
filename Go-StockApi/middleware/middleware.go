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

// response struct for API responses
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message.omitempty"`
}

//connectDatabase establishes a connecttion to the PostGres database
/*
Returns :
	- *sql.DB: A pointer to the sql DataBase connection
*/
func connectDatabase() *sql.DB {

	DB_HOST := config.GetEnv("DB_HOST")
	DB_PORT := config.GetEnv("DB_PORT")
	DB_USER := config.GetEnv("DB_USER")
	DB_PASSWORD := config.GetEnv("DB_PASSWORD")
	DB_NAME := config.GetEnv("DB_NAME")
	dsn_info := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", dsn_info)
	if err != nil {
		log.Fatal("❌ Database connection failed: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("❌ Database ping failed: ", err)
	}

	fmt.Println("✅ Successfully connected to the database")
	return db
}

//CreateStock handles the creation of a new Stock
/*
Args:
	- w	(http.ResponseWriter): Response Writer to send responses back to the client
	-r (*http.Request): HTTP request containing the stock data in the request body

	Returns:
	- JSON response with inserted stock ID and success message
*/
func CreateStock(w http.ResponseWriter, r *http.Request) {

	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("❌ Unable to decode the request body: %v", err)
	}

	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "✅ Stock Created Successfully",
	}

	json.NewEncoder(w).Encode(res)
}

//Get Stock fetches a stock by its ID
/*
	Args:
		- w (http.ResponseWriter): Response Writer
		-r (*http.Request): HTTP request containing stock ID as a URL parameter
	Returns:
		- Json response with stock details
*/
func GetStock(w http.ResponseWriter, r *http.Request) {
	parmas := mux.Vars(r)

	id, err := strconv.Atoi(parmas["id"])
	if err != nil {
		log.Fatalf("❌ Unable to convert String to INT, %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("❌ Unable to get strock %v", err)
	}
	json.NewEncoder(w).Encode((stock))
}

// GetAllStock fetches all stocks from the database.
/*
   Args:
   - w (http.ResponseWriter): Response writer.
   - r (*http.Request): HTTP request.

   Returns:
   - JSON response containing all stock records.
*/
func GetAllStock(w http.ResponseWriter, r *http.Request) {

	stocks, err := getAllstocks()

	if err != nil {
		log.Fatalf("❌ Unable to fetch all stocks. %v ", err)
	}
	json.NewEncoder(w).Encode(stocks)

}

// UpdateStock updates an existing stock record.
/*
   Args:
   - w (http.ResponseWriter): Response writer.
   - r (*http.Request): HTTP request containing stock ID and updated stock details.

   Returns:
   - JSON response with update status.
*/
func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("❌ Unable to convert string to into int. %v", err)
	}
	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("❌ Unable to decode Request Body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("✅Stock Updated successfully. Total row/ record affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

// DeleteStock removes a stock record from the database.
/*
   Args:
   - w (http.ResponseWriter): Response writer.
   - r (*http.Request): HTTP request containing stock ID.

   Returns:
   - JSON response with delete status.
*/
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("❌ Unable to convert string to integer %v", err)
	}
	deleterRows := deleteStock(int64(id))

	msg := fmt.Sprintf("✅ Successfully updated stock . Total Rows/record Affected %v", deleterRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

//---------------------- DATABASE OPERATIONS ----------------------

// insertStock inserts a new stock record into the database.
/*
   Args:
   - stock (models.Stock): Stock object containing name, price, and company.

   Returns:
   - int64: ID of the newly inserted stock.
*/
func insertStock(stock models.Stock) int64 {

	db := connectDatabase()
	defer db.Close()

	sqlStatement := "INSERT INTO STOCK  (name,price,company) VALUES ($1 ,$2, $3)RETURNING stockid"
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("❌ Unable to insert stock %v", err)
	}
	fmt.Printf("✅ Inserted Single Record. %v", id)
	return id
}

// getStock retrieves a single stock record by ID.
/*
   Args:
   - id (int64): Stock ID.

   Returns:
   - models.Stock: Stock object.
   - error: Error if any issue occurs.
*/
func getStock(id int64) (models.Stock, error) {
	db := connectDatabase()
	defer db.Close()

	var stock models.Stock
	sqlStament := " SELECT * FROM stock WHERE stockid = $1"
	row := db.QueryRow(sqlStament, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No Stock Fonnd")
			return stock, nil
		}
		log.Fatalf("❌ Unable to scan row: %v", err)

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
