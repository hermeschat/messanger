/*
type User struct {
	Id int
	Name string
}
func (u *User) Scan(row *sql.Rows) error {
	if !row.Next() {
		return errors.New("no result in query response")
	}
	return row.Scan(&u.Id, &u.Name)
}
func (u *User) Insert(db *sql.DB) (sql.Result, error) {
	return db.Exec(InsertQuery("users", "id", "name"), u.Id, u.Name)
}
func (u *User) Update(db *sql.DB) (sql.Result, error) {
	return db.Exec(UpdateQuery("users", "id=?", "name"),u.Id, u.Name)
}
func (u *User) Delete(db *sql.DB) (sql.Result, error) {
	return db.Exec(DeleteQuery("users", "id=?"), u.Id)
}

func main() {
	var db *sql.DB
	u := &User{Id:1}
	Query(db, u, "WHERE id=?", u.Id)
	u.Name = "amirreza"
	u.Update(db)
	u.Delete(db)
	u2 := &User{Name:"Asghar"}
	u2.Insert(db)
}

*/

package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type Selectable interface {
	Scan(*sql.Rows) error
}

type Insertable interface {
	Insert() (string, []interface{})
}

type Updateable interface {
	Update() (string, []interface{})
}

type Deleteable interface {
	Delete() (string, []interface{})
}

type SqlEntity interface {
	Selectable
	Insertable
	Updateable
	Deleteable
}

//Query
func Query(db *sql.DB, rs Selectable, query string, args ...interface{}) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rs.Scan(rows)
		if err != nil {
			return err
		}
	}
	return nil
}
func Insert(db *sql.DB, i Insertable) (sql.Result, error) {
	stmt, values := i.Insert()
	return db.Exec(stmt, values...)
}
//InsertQuery is a helper method to generate insert query for given table and columns.
func InsertQuery(table string, columns ...string) string {
	var questions []string
	for i:=0;i<len(columns);i++ {
		questions = append(questions, "?")
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(questions, ", "))
}
func Update(db *sql.DB, u Updateable) (sql.Result, error) {
	stmt, values := u.Update()
	return db.Exec(stmt, values...)
}
//UpdateQuery is a helper method to generate an update query for given table and columns.
func UpdateQuery(table string, where string, columns ...string) string {
	var colPair []string
	for _, c := range columns {
		colPair = append(colPair, fmt.Sprintf("%s=?", c))
	}
	return fmt.Sprintf("UPDATE %s WHERE %s SET %s", table, where, strings.Join(colPair, ", "))
}
func Delete(db *sql.DB, d Deleteable) (sql.Result, error) {
	stmt, values := d.Delete()
	return db.Exec(stmt, values...)
}
//DeleteQuery is a helper method to generate a delete query for given table.
func DeleteQuery(table string, where string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
}
