package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

func main() {
	conn, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	s := server{db: db{conn: conn}}
	mux := http.NewServeMux()
	mux.HandleFunc("/", website)
	mux.HandleFunc("/FindInvite", s.findInvite)
	mux.HandleFunc("/GetInvite", s.getInvite)
	mux.HandleFunc("/RSVP", s.rsvp)
	mux.HandleFunc("/UpdateEmail", s.updateEmail)

	log.Fatal(http.ListenAndServe(os.Args[1], cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
	}).Handler(mux)))
}

func website(w http.ResponseWriter, r *http.Request) {
	var homeTemplate = template.New("Home")
	if html, err := os.ReadFile("website.html"); err != nil {
		panic(err.Error())
	} else {
		if _, err := homeTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			homeTemplate.Execute(w, nil)
		}
	}
}

type server struct {
	db db
}

func (s *server) findInvite(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	invites, err := s.db.getFuzzyInvite(name)
	if err != nil {
		retErr(w, err)
		return
	}

	b, err := json.MarshalIndent(invites, "", "\t")
	if err != nil {
		retErr(w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func (s *server) getInvite(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		retErr(w, err)
		return
	}

	persons, err := s.db.getInvite(id)
	if err != nil {
		retErr(w, err)
		return
	}

	b, err := json.MarshalIndent(persons, "", "\t")
	if err != nil {
		retErr(w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func (s *server) rsvp(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		retErr(w, err)
		return
	}

	var p Person
	err = json.Unmarshal(b, &p)
	if err != nil {
		retErr(w, err)
		return
	}

	err = s.db.rsvp(p)
	if err != nil {
		retErr(w, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("done"))
}

func (s *server) updateEmail(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		retErr(w, err)
		return
	}

	var body struct {
		Id    int
		Email string
	}
	err = json.Unmarshal(b, &body)
	if err != nil {
		retErr(w, err)
		return
	}

	err = s.db.updateEmail(body.Id, body.Email)
	if err != nil {
		retErr(w, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("done"))
}

func retErr(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
