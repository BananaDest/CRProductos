package db

import (
	"crproductos/models"
	"crproductos/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectToPostgres() *sql.DB {
	err := godotenv.Load("../config.development.env")
	if err != nil {
		log.Fatalf("Error loading .env files: %v", err)
	}
	dbhost := os.Getenv("DB_HOST")
	dbportStr := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbport, err := strconv.Atoi(dbportStr)
	if err != nil {
		log.Fatalf("error converting: %v", err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", psqlInfo)
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

func GetAllProducts(DB *sql.DB) ([]models.ProductResponse, error) {
	rows, err := DB.Query("select id,\"name\",quantity,unit,stores from product")
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

func GetProductById(DB *sql.DB, id string) (models.Product, error) {
	var product models.Product
	if err := DB.QueryRow("select * from product where product.id = $1", id).Scan(&product.Id, &product.Name, &product.Quantity, &product.Unit, &product.Stores); err != nil {
		log.Println("failed to scan: ", err)
		return product, err
	}
	log.Printf("Item found: %+v\n", product.ToJSON())
	return product, nil
}

func CreateProduct(DB *sql.DB, product models.ProductResponse) (models.ProductResponse, error) {
	tx, err := DB.Begin()
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

func DeleteProduct(DB *sql.DB, id string) error {
	tx, err := DB.Begin()
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

func UpdateProduct(DB *sql.DB, id string, product models.ProductResponse) (models.ProductResponse, error) {
	tx, err := DB.Begin()
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

func PatchProduct(DB *sql.DB, id string, product models.ProductResponse) (models.Product, error) {
	tx, err := DB.Begin()
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
func PatchStore(DB *sql.DB, id string, jsonStore []byte) (models.Product, error) {

	var updatedProduct models.Product
	tx, err := DB.Begin()
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
