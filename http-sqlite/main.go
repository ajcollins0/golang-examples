package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Server - main server struct
type Server struct {
	Serv *http.ServeMux
	Db   *sql.DB
}

// CreateTables - creates tables so db is at latest schema
func (s *Server) createTables() {
	sqlStmt := `
	create table if not exists hosts (
		host     text primary key
	);`
	_, err := s.Db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

// OpenDB - Opens and creates DB
func (s *Server) OpenDB(dbPath string) {
	var err error
	s.Db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	s.createTables()
}

// CloseDB - Close DB
func (s *Server) CloseDB() {
	s.Db.Close()
}

// CreateHost - add host information to DB
func CreateHost(s *Server) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		host := r.URL.Query()["target"][0]
		stmt, err := s.Db.Prepare("INSERT INTO hosts(host) values(?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err := stmt.Exec(host)
		fmt.Println(res)
		if err != nil {
			log.Fatal(err)
		}
	}
	return http.HandlerFunc(fn)
}

func main() {

	dbSTr := "./test.db"

	s := &Server{new(http.ServeMux), new(sql.DB)}

	s.OpenDB(dbSTr)
	defer s.CloseDB()

	s.Serv.HandleFunc("/addhost", CreateHost(s))

	log.Fatal(http.ListenAndServe(":8080", s.Serv))
}
