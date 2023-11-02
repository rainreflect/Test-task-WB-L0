package server

import (
	"html/template"
	"log"
	"net/http"
	"rec/internal/pkg/postgresql"
	caching "rec/internal/pkg/postgresql/caching"
	"rec/internal/pkg/postgresql/model"

	"github.com/gorilla/mux"
)

type Server struct {
	Serv     *http.Server
	storage  caching.CacheService
	database *postgresql.DBService
	Address  string
}

func NewServer(s caching.CacheService, a string, db *postgresql.DBService) *Server {
	srv := Server{
		storage:  s,
		Address:  a,
		database: db,
	}
	return &srv
}

func (s *Server) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/orders/{o_id}", s.ordersHandler)
	s.Serv = &http.Server{Addr: s.Address, Handler: router}
	log.Println("SERVER HOSTING")
	err := s.Serv.ListenAndServe()

	return err
}

func (s *Server) Stop() error {
	log.Println("SERVER CLOSING")
	return s.Serv.Close()
}

func (s *Server) ordersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["o_id"]

	od, err := s.database.OrderById(id)
	if err != nil {
		log.Println("error while taking orders by id")
	}
	if od.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		parsedTemplate, err1 := template.ParseFiles("internal/pkg/server/UI/no.html")
		if err1 != nil {
			log.Println(err1)
		}
		err := parsedTemplate.Execute(w, struct{ Id string }{Id: id})
		if err != nil {
			_, err = w.Write([]byte("no data with id " + id))
			if err != nil {
				return
			}
			log.Println("Error occurred while executing the template : ", id)
			return
		}
	} else if od.ID >= "1" {
		orderItem := model.OrderItem{
			ID:   id,
			Data: od.Data,
		}
		parsedTemplate, _ := template.ParseFiles("internal/pkg/server/UI/index.html")
		err := parsedTemplate.Execute(w, orderItem)
		if err != nil {
			w.Write([]byte("error while executing template"))
			log.Println("Error occurred while executing the template : ", orderItem)
			return
		}
	}
}
