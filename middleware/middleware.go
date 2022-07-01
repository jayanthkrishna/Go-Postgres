package middleware

import (
	"Go-Postgres/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Name    string `json:"stockname omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error Loading .env File")

	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db

}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("Unable to decode the request body. ", err)
	}

	stockID, stockName := insertStock(stock)

	res := Response{
		ID:      stockID,
		Name:    stockName,
		Message: "Stock Created Successfully",
	}

	json.NewEncoder(w).Encode(res)

}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatal("Unable to convert string into int", err)
	}
	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatal("Unable to get stock ", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {

	stocks, err := getAllStocks()

	if err != nil {
		log.Fatal("Unable to get all the stocks", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatal("Unable to convert string into int", err)
	}
	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to deocde the request body, %v\n", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock Updated Successfully. No of rows affected are %v", updatedRows)

	res := Response{
		ID:      int64(id),
		Name:    stock.Name,
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string into int %v", err)
	}
	stock, _ := getStock(int64(id))
	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock Deleted Successfully. Number of rows affected are : %v", deletedRows)

	res := Response{
		ID:      stock.StockID,
		Name:    stock.Name,
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()

	defer db.Close()

	var stock models.Stock

	sqlstatement := `SELECT * FROM STOCKS WHERE STOCKID = $1`

	row := db.QueryRow(sqlstatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatal("Unable to read the row")

	}

	return stock, err

}

func insertStock(stock models.Stock) (int64, string) {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO STOCKS(NAME,PRICE,COMPANY) VALUES($1,$2,$3) RETURNING STOCKID`
	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to Execute the Query: %v", err)
	}

	fmt.Printf("Inserted a seingle record with ID : %d", id)

	return id, stock.Name

}

func getAllStocks() ([]models.Stock, error) {
	db := createConnection()

	defer db.Close()

	var stocks []models.Stock

	sqlstatement := `SELECT * FROM STOCKS`

	rows, err := db.Query(sqlstatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var stock models.Stock

		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan the row %v", err)
		}
		stocks = append(stocks, stock)

	}

	return stocks, err

}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()

	defer db.Close()

	sqlstatement := `UPDATE STOCKS SET NAME=$2,PRICE=$3,COMPANY=$4 WHERE STOCKID=$1`

	res, err := db.Exec(sqlstatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Unable to excute the query %v", err)
	}

	rowsAffected, _ := res.RowsAffected()

	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := createConnection()

	defer db.Close()

	sqlstatement := `DELETE FROM STOCKS WHERE STOCKID=$1`

	res, err := db.Exec(sqlstatement, id)

	if err != nil {
		log.Fatalf("Unable to excute the query %v", err)
	}

	rowsAffected, _ := res.RowsAffected()

	return rowsAffected

}
