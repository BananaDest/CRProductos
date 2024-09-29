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
	"reflect"
	"strings"
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
		r.Get("/", server.GetAllProducts)
		r.Get("/{id}", server.GetProductById)
		r.Post("/", server.CreateProduct)
		r.Put("/{id}", server.UpdateProduct)
		r.Delete("/{id}", server.DeleteProduct)
		r.Patch("/{id}", server.PatchProduct)
	})
	http.ListenAndServe(":8080", r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
func (s *Server) GetAllProducts(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) GetProductById(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var product Product
	if err := s.DB.QueryRow("select * from product where product.id = $1", id).Scan(&product.id, &product.name, &product.quantity, &product.unit, &product.stores); err != nil {
		log.Println("failed to scan: ", err)

		return
	}
	log.Printf("Item found: %+v\n", product.ToJSON())

	render.JSON(w, r, product.ToJSON())

}
func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Printf("producto %+v\n", product)
	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("INSERT INTO public.product (\"name\", quantity, unit, stores) VALUES($1, $2, $3, $4);", product.Name, product.Quantity, product.Unit, product.Stores)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error during insert: ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commiting: ", err)
	}
	render.JSON(w, r, product)
}
func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")

	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("DELETE FROM public.product WHERE id=$1;", id)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error during delete: ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commiting: ", err)
	}
	w.Write([]byte("Delete successful"))
}

func (s *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var product ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("UPDATE public.product SET \"name\"=$1, quantity=$2, unit=$3, stores=$4 WHERE id=$5;", product.Name, product.Quantity, product.Unit, product.Stores, id)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error during update: ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commiting: ", err)
	}
	render.JSON(w, r, product)

}
func (s *Server) PatchProduct(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	var product ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	var updateClauses []string
	var args []interface{}
	argIndex := 1
	v := reflect.ValueOf(product)
	t := reflect.TypeOf(product)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Check if the field is a pointer
		if field.Kind() == reflect.Ptr {
			// Only check if the pointer is not nil
			if !field.IsNil() {
				columnName := fieldToColumnName(t.Field(i)) // Get the corresponding column name from the struct tag
				fieldValue := field.Elem().Interface()      // Dereference the pointer to get the actual value

				// Add to update clauses
				updateClauses = append(updateClauses, fmt.Sprintf("%s = $%d", columnName, argIndex))
				args = append(args, fieldValue)
				argIndex++
			}
		} else {
			// Handle non-pointer fields directly
			if !field.IsZero() { // Check if the value is the zero value for its type
				columnName := fieldToColumnName(t.Field(i))
				fieldValue := field.Interface() // Get the actual value

				// Add to update clauses
				updateClauses = append(updateClauses, fmt.Sprintf("%s = $%d", columnName, argIndex))
				args = append(args, fieldValue)
				argIndex++
			}
		}

	}
	if len(updateClauses) == 0 {
		http.Error(w, "no fields for update", http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf("Update product set %s where id=$%d", strings.Join(updateClauses, ", "), argIndex)
	args = append(args, id)
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}
	var updatedProduct Product
	tx.QueryRow("select * from product where product.id = $1", id).Scan(&updatedProduct.id, &updatedProduct.name, &updatedProduct.quantity, &updatedProduct.unit, &updatedProduct.stores)
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commiting: ", err)
	}
	render.JSON(w, r, updatedProduct.ToJSON())

}
func fieldToColumnName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}
	return strings.Split(jsonTag, ",")[0]
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
