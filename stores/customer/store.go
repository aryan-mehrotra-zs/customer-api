package customer

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/amehrotra/customer-api/model"
)

type store struct {
	db *sql.DB
}

func New(db *sql.DB) store {
	return store{db: db}
}

func (s store) Get(id int) (model.Customer, error) {
	var c model.Customer

	err := s.db.QueryRow("SELECT * FROM customers WHERE ID = ?", id).
		Scan(&c.ID, &c.Name, &c.Address, &c.PhoneNo)

	return c, err
}

func (s store) Create(c model.Customer) (model.Customer, error) {
	result, err := s.db.Exec("INSERT INTO customers (id,name,address,phone_no) VALUES (?,?,?,?)", c.ID, c.Name, c.Address, c.PhoneNo)
	if err != nil {
		return model.Customer{}, err
	}

	data, err := result.LastInsertId()
	if err != nil {
		return model.Customer{}, err
	}

	return s.Get(int(data))
}

func (s store) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM customers WHERE id = ?;", id)

	return err
}

func (s store) Update(c model.Customer) error {
	query := createPutQuery(c.ID, c)
	if query == "" {
		return nil
	}

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func createPutQuery(id int, c model.Customer) string {
	var q []string

	if c.Name != "" {
		q = append(q, " name = \""+c.Name+"\"")
	}

	if c.Address != "" {
		q = append(q, " address = \""+c.Address+"\"")
	}

	if c.PhoneNo != 0 {
		q = append(q, " phone_no = "+strconv.Itoa(c.PhoneNo))
	}

	if q == nil {
		return ""
	}

	query := "UPDATE customers SET" + strings.Join(q, ",") + " WHERE id=" + strconv.Itoa(id) + ";"

	return query
}
