package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"html/template"
	"l0/internal/database"
	"l0/internal/model"
	"net/http"
)

type config struct {
	Host             string
	Port             string
	SubscribeSubject string
}

func newConfig() *config {
	config := config{Host: "localhost", Port: ":8000", SubscribeSubject: "addNewOrder"}
	return &config
}

type Server struct {
	cache  map[string]model.Order
	db     *database.Database
	config *config
	router *chi.Mux
	sc     stan.Conn
	sub    stan.Subscription
}

func setConfigs() (*database.Database, *config) {
	db := database.SetConfig()
	config := newConfig()
	return db, config
}

func NewServer() *Server {
	db, config := setConfigs()
	return &Server{
		db:     db,
		cache:  make(map[string]model.Order),
		config: config,
		router: chi.NewRouter(),
	}
}

func (s *Server) connectToStream() error {
	sc, err := stan.Connect("test-cluster", "subscriber", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return err
	}
	sub, err := sc.Subscribe(s.config.SubscribeSubject, s.handleRequest)
	if err != nil {
		return err
	}
	s.sc, s.sub = sc, sub
	return nil
}

func (s *Server) Up() error {
	s.db.Open()
	if err := s.setCache(); err != nil {
		return err
	}
	if err := s.connectToStream(); err != nil {
		return err
	}
	s.startRouting()
	return nil
}

func (s *Server) Down() {
	logrus.Info("Server is down")
	s.db.Close()
	err := s.sub.Unsubscribe()
	if err != nil {
		return
	}
	err = s.sc.Close()
	if err != nil {
		return
	}
}

func (s *Server) handleRequest(m *stan.Msg) {
	data := model.Order{}
	err := json.Unmarshal(m.Data, &data)
	if err != nil {
		return
	}
	if ok := s.addToCache(data); ok {
		logrus.Info("Add to cache")
		err := s.db.AddOrder(data)
		if err != nil {
			return
		}
	}
}

func (s *Server) setCache() error {
	orders := make([]model.Order, 0)
	err := s.db.DB.Model(&orders).Select()
	if err != nil {
		return err
	}
	for _, order := range orders {
		s.cache[order.OrderUid] = order
	}
	return nil
}

func (s *Server) addToCache(data model.Order) bool {
	_, ok := s.cache[data.OrderUid]
	if ok {
		return false
	}
	s.cache[data.OrderUid] = data
	for key := range s.cache {
		fmt.Printf("%s ", key)
	}
	fmt.Println()
	return true
}

func (s *Server) startRouting() {
	s.router.Use(middleware.Logger)
	s.router.Get("/", s.WelcomeHandler)
	s.router.Get("/order/{order_uid}", s.handleGetId)
	address := fmt.Sprintf("%s%s", s.config.Host, s.config.Port)
	logrus.Info("Server is up")
	err := http.ListenAndServe(address, s.router)
	if err != nil {
		return
	}
}

func (s *Server) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = t.ExecuteTemplate(w, "home_page.html", nil)
	if err != nil {
		return
	}
}

func (s *Server) handleGetId(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "order_uid")
	str := fmt.Sprintf("Order with id %s has been added", id)
	logrus.Info(str)
	data, ok := s.cache[id]
	if !ok {
		writer.Write([]byte("Nothing was found. Please, try another id."))
		return
	}
	tmpl, _ := template.ParseFiles("templates/home_page.html")
	err := tmpl.Execute(writer, data)
	if err != nil {
		return
	}
}
