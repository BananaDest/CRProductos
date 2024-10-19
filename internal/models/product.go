package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Product struct {
	Id       int
	Name     sql.NullString
	Quantity sql.NullFloat64
	Unit     sql.NullString
	Stores   *Stores
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
	if p.Name.Valid {
		name = &p.Name.String
	}
	var quantity *float64
	if p.Quantity.Valid {
		quantity = &p.Quantity.Float64
	}
	var unit *string
	if p.Unit.Valid {
		unit = &p.Unit.String
	}

	return ProductResponse{
		Id:       p.Id,
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
		Stores:   p.Stores,
	}
}

func (p *ProductResponse) ToString() string {
	storesString := "{"
	for key, storeValue := range *p.Stores {
		storesString = storesString + fmt.Sprintf("%+v: %+v", key, storeValue) + ", "
	}
	if len(storesString) >= 3 {
		storesString = storesString[:len(storesString)-2]
	}
	storesString += "}"
	productString := fmt.Sprintf("{Id: %d: Name: %s, Quantity: %f, Unit: %s, Stores: %s}", p.Id, *p.Name, *p.Quantity, *p.Unit, storesString)
	return productString
}
