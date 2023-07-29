package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type studentinfo struct {
	Sid   string `json:"sid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(appmoc.com:3306)/studentinfo?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	defer db.Close()
	ss := []studentinfo{}
	s := studentinfo{}
	rows, err := db.Query("select * from student")
	if err != nil {
		fmt.Fprintf(w, "this is "+err.Error())
	} else {

		for rows.Next() {
			rows.Scan(&s.Sid, &s.Name, &s.Email)
			ss = append(ss, s)
		}
		json.NewEncoder(w).Encode(ss)
	}

}

func addStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	defer db.Close()
	s := studentinfo{}
	json.NewDecoder(r.Body).Decode(&s)
	sid, _ := strconv.Atoi(s.Sid)
	result, err := db.Exec("insert into student(sid,name,email) values(?,?,?)", sid, s.Name, s.Email)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())

	} else {
		_, err := result.LastInsertId()
		if err != nil {
			json.NewEncoder(w).Encode("{error:Record not inserted}")
		} else {
			json.NewEncoder(w).Encode(s)
		}
	}

}

func updateStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	defer db.Close()
	s := studentinfo{}
	json.NewDecoder(r.Body).Decode(&s)
	vars := mux.Vars(r)

	sid, _ := strconv.Atoi(vars["sid"])
	result, err := db.Exec("update student set name=?, email=?  where sid=?", s.Name, s.Email, sid)
	if err != nil {

		fmt.Fprintf(w, ""+err.Error())

	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{error: Record not updated}")

		} else {
			json.NewEncoder(w).Encode(s)
		}
	}

	fmt.Fprintf(w, "UPDATE STUDENTS")
}

func DeleteStudents(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	defer db.Close()

	fmt.Fprintf(w, "DELETE STUDENTS")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/studentslist", getStudents).Methods("GET")
	r.HandleFunc("/students", addStudents).Methods("POST")
	r.HandleFunc("/students/{sid}", updateStudents).Methods("PATCH")
	r.HandleFunc("/students/{sid}", DeleteStudents).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}
