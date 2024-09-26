package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "crproducts"
)

type Product struct {
	id       int
	name     sql.NullString
	quantity sql.NullFloat64
	unit     sql.NullString
	stores   *Stores
}

type Stores map[string]float64

func (s Stores) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Stores) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to byte failed")
	}

	return json.Unmarshal(b, s)
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db := connectToPostgres(psqlInfo)
	defer db.Close()
	rows, err := db.Query("select id,\"name\",quantity,unit,stores from product")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		//TODO: generate structure and add rows
		//TODO: test driven development
		var product Product
		if err := rows.Scan(&product.id, &product.name, &product.quantity, &product.unit, &product.stores); err != nil {
			log.Fatal(err)
		}
		fmt.Println(product.stores)
		products = append(products, product)
	}
	fmt.Print(products)
	r := chi.NewRouter()
	r.Get("/", rootHandler)
	http.ListenAndServe(":8080", r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func connectToPostgres(connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect ", err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping postgres", err)
	}
	fmt.Println("connected to postgres")
	return db
}
