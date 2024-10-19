package repository

import (
	"crproductos/internal/models"
	"crproductos/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"strings"
)

type userRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &userRepository{db: db}
}

// GetAllProducts: Receives the r.db struct instance and returns
// either a list of all products found, or the corresponding error.
// A successful GetAllProducts call will return err == nil
func (r *userRepository) GetAllProducts() ([]models.ProductResponse, error) {
	rows, err := r.db.Query("select id,\"name\",quantity,unit,stores from product")
	if err != nil {
		log.Println("Failed to query product: ", err)
		return nil, err
	}
	defer rows.Close()
	var products []models.ProductResponse
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Quantity, &product.Unit, &product.Stores); err != nil {
			log.Println("failed to scan: ", err)
			return nil, err
		}
		log.Printf("Item found: %+v\n", product.ToJSON())
		products = append(products, product.ToJSON())
	}
	if err = rows.Err(); err != nil {
		log.Println("row iteration error:", err)
		return nil, err
	}
	return products, nil
}

func (r *userRepository) GetProductById(id string) (models.Product, error) {
	var product models.Product
	if err := r.db.QueryRow("select * from product where product.id = $1", id).Scan(&product.Id, &product.Name, &product.Quantity, &product.Unit, &product.Stores); err != nil {
		log.Println("failed to scan: ", err)
		return product, err
	}
	log.Printf("Item found: %+v\n", product.ToJSON())
	return product, nil
}

func (r *userRepository) CreateProduct(product models.ProductResponse) (models.ProductResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	err = tx.QueryRow("INSERT INTO public.product (\"name\", quantity, unit, stores) VALUES($1, $2, $3, $4) returning id;", product.Name, product.Quantity, product.Unit, product.Stores).Scan(&product.Id)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error during insert: ", err)
		return product, err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commiting: ", err)
	}
	return product, nil
}

func (r *userRepository) DeleteProduct(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("DELETE FROM public.product WHERE id=$1;", id)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error during delete: ", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commiting: ", err)
	}
	return nil
}

func (r *userRepository) UpdateProduct(id string, product models.ProductResponse) (models.ProductResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("UPDATE public.product SET \"name\"=$1, quantity=$2, unit=$3, stores=$4 WHERE id=$5;", product.Name, product.Quantity, product.Unit, product.Stores, id)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error during update: ", err)
		return product, err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commiting: ", err)
	}
	return product, nil
}

func (r *userRepository) PatchProduct(id string, product models.ProductResponse) (models.Product, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	var updateClauses []string
	var args []interface{}
	argIndex := 1
	v := reflect.ValueOf(product)
	t := reflect.TypeOf(product)

	var updatedProduct models.Product
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
		fmt.Println("No fields to update")
		return updatedProduct, errors.New("No fields for update")
	}
	query := fmt.Sprintf("Update product set %s where id=$%d", strings.Join(updateClauses, ", "), argIndex)
	args = append(args, id)
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return updatedProduct, err
	}
	tx.QueryRow("select * from product where product.id = $1", id).Scan(&updatedProduct.Id, &updatedProduct.Name, &updatedProduct.Quantity, &updatedProduct.Unit, &updatedProduct.Stores)
	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error commiting: %v", err)
		return updatedProduct, err
	}
	return updatedProduct, nil
}
func (r *userRepository) PatchStore(id string, jsonStore []byte) (models.Product, error) {

	var updatedProduct models.Product
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	query := "Update product set stores = stores || $1::jsonb where id=$2"

	_, err = tx.Exec(query, string(jsonStore), id)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to Patch: %v\n", err)
		return updatedProduct, err
	}
	tx.QueryRow("select * from product where product.id = $1", id).Scan(&updatedProduct.Id, &updatedProduct.Name, &updatedProduct.Quantity, &updatedProduct.Unit, &updatedProduct.Stores)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to Patch: %v\n", err)
		return updatedProduct, err
	}
	return updatedProduct, nil
}
