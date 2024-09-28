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
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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

type ProductResponse struct {
	Id       int      `json:"id"`
	Name     *string  `json:"name"`
	Quantity *float64 `json:"quantity"`
	Unit     *string  `json:"unit"`
	Stores   *Stores  `json:"stores"`
}

func (p *Product) ToJSON() ProductResponse {
	var name *string
	if p.name.Valid {
		name = &p.name.String
	}
	var quantity *float64
	if p.quantity.Valid {
		quantity = &p.quantity.Float64
	}
	var unit *string
	if p.unit.Valid {
		unit = &p.unit.String
	}

	return ProductResponse{
		Id:       p.id,
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
		Stores:   p.stores,
	}
}

type Server struct {
	DB *sql.DB
}

func newServer(db *sql.DB) *Server {
	return &Server{DB: db}
}

func (s *Server) getAllProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query("select id,\"name\",quantity,unit,stores from product")
	if err != nil {
		http.Error(w, "Failed to query product", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []ProductResponse
	for rows.Next() {
		//TODO: test driven development
		var product Product
		if err := rows.Scan(&product.id, &product.name, &product.quantity, &product.unit, &product.stores); err != nil {
			log.Println("failed to scan: ", err)
			http.Error(w, "Failed to scan product", http.StatusInternalServerError)
			return
		}
		log.Printf("Item found: %+v\n", product.ToJSON())
		products = append(products, product.ToJSON())
	}
	if err = rows.Err(); err != nil {
		log.Println("row iteration error:", err)
		http.Error(w, "Failed to iterate rows", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, products)

}
func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db := connectToPostgres(psqlInfo)
	defer db.Close()
	server := newServer(db)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", rootHandler)
	r.Route("/products", func(r chi.Router) {
		r.Get("/", server.getAllProducts)
		// r.Get("/{id}", getProductByIdHandler)
		// r.Post("/", createProductHandler)
		// r.Put("/{id}", updateProductHandler)
		// r.Delete("/{id}", deleteProductHandler)
	})
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
