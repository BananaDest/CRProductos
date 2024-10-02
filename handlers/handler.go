package handlers

import (
	"crproductos/models"
	"crproductos/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type Server struct {
	DB *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{DB: db}
}
func (s *Server) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query("select id,\"name\",quantity,unit,stores from product")
	if err != nil {
		http.Error(w, "Failed to query product", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []models.ProductResponse
	for rows.Next() {
		//TODO: test driven development
		var product models.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Quantity, &product.Unit, &product.Stores); err != nil {
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
	var product models.Product
	if err := s.DB.QueryRow("select * from product where product.id = $1", id).Scan(&product.Id, &product.Name, &product.Quantity, &product.Unit, &product.Stores); err != nil {
		log.Println("failed to scan: ", err)

		return
	}
	log.Printf("Item found: %+v\n", product.ToJSON())

	render.JSON(w, r, product.ToJSON())

}
func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.ProductResponse
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
	var product models.ProductResponse
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
	var product models.ProductResponse
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
				columnName := utils.FieldToColumnName(t.Field(i)) // Get the corresponding column name from the struct tag
				fieldValue := field.Elem().Interface()            // Dereference the pointer to get the actual value

				// Add to update clauses
				updateClauses = append(updateClauses, fmt.Sprintf("%s = $%d", columnName, argIndex))
				args = append(args, fieldValue)
				argIndex++
			}
		} else {
			// Handle non-pointer fields directly
			if !field.IsZero() { // Check if the value is the zero value for its type
				columnName := utils.FieldToColumnName(t.Field(i))
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
	var updatedProduct models.Product
	tx.QueryRow("select * from product where product.id = $1", id).Scan(&updatedProduct.Id, &updatedProduct.Name, &updatedProduct.Quantity, &updatedProduct.Unit, &updatedProduct.Stores)
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commiting: ", err)
	}
	render.JSON(w, r, updatedProduct.ToJSON())

}
