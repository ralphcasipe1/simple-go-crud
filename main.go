package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	ContactID int
	Name      string
	Company   string
	Phone     string
	Address   string
	Email     string
	Photo     string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "supersun"
	dbName := "contact_management"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT contact_id, name, company, phone, address, email FROM contacts ORDER BY contact_id DESC")
	if err != nil {
		panic(err.Error())
	}

	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var contactID int
		var name, company, phone, address, email string

		err = selDB.Scan(&contactID, &name, &company, &phone, &address, &email)

		if err != nil {
			panic(err.Error())
		}

		emp.ContactID = contactID
		emp.Name = name
		emp.Company = company
		emp.Phone = phone
		emp.Address = address
		emp.Email = email
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nID := r.URL.Query().Get("contact_id")
	selDB, err := db.Query("SELECT * FROM contacts WHERE contact_id=?", nID)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var contactID int
		var name, company string
		err = selDB.Scan(&contactID, &name, &company)
		if err != nil {
			panic(err.Error())
		}
		emp.ContactID = contactID
		emp.Name = name
		emp.Company = company
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		company := r.FormValue("company")
		insForm, err := db.Prepare("INSERT INTO contacts (name, company) VALUES (?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, company)
		log.Println("INSERT: Name: " + name + " | Company: " + company)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		company := r.FormValue("company")
		contactID := r.FormValue("contact_id")
		insForm, err := db.Prepare("UPDATE contacts SET name=?, company=? WHERE contact_id=?")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, company, contactID)
		log.Println("UPDATE: Name: " + name + " | Company: " + company)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	contact := r.URL.Query().Get("contact_id")
	delForm, err := db.Prepare("DELETE FROM contacts WHERE contact_id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(contact)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", New)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
