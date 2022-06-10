package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id      int
	Name    string
	Contact string
	Address string
	Dob     string
	Gender  string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "godata"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, contact, address, dob, gender string
		err = selDB.Scan(&id, &name, &contact, &address, &dob, &gender)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Contact = contact
		emp.Address = address
		emp.Dob = dob
		emp.Gender = gender
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, contact, address, dob, gender string
		err = selDB.Scan(&id, &name, &contact, &address, &dob, &gender)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Contact = contact
		emp.Address = address
		emp.Dob = dob
		emp.Gender = gender
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, contact, address, dob, gender string
		err = selDB.Scan(&id, &name, &contact, &address, &dob, &gender)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Contact = contact
		emp.Address = address
		emp.Dob = dob
		emp.Gender = gender
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		contact := r.FormValue("contact")
		address := r.FormValue("address")
		dob := r.FormValue("dob")
		gender := r.FormValue("gender")
		insForm, err := db.Prepare("INSERT INTO Employee(name, contact, address, dob, gender) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, contact, address, dob, gender)
		log.Println(insForm)

		log.Println("INSERT: Name: " + name + " | Contact: " + contact + " | Address: " + address + " | Dob: " + dob + " | Gender: " + gender)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		contact := r.FormValue("contact")
		address := r.FormValue("address")
		dob := r.FormValue("dob")
		gender := r.FormValue("gender")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Employee SET name=?, contact=?, address=?, dob=?, gender=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, contact, address, dob, gender, id)
		log.Println("UPDATE: Name: " + name + " | Contact: " + contact + " | Address: " + address + " | Dob: " + dob + " | Gender: " + gender)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	log.Print("DELETE is qworking")
	db := dbConn()

	emp := r.URL.Query().Get("id")

	log.Println(emp)

	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8282")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/sally", Delete)

	http.ListenAndServe(":8282", nil)
}
