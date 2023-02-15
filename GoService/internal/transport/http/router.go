package http

import (
	"context"

	"example.com/GoService/pkg/models"

	"github.com/gorilla/mux"
)

type service interface {
	GetDataByID(ctx context.Context, uID string) (order models.Order, err error)
}

type Router struct {
	svc service
}

func (rt *Router) InitRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/data/{id:[a-z;A-Z;0-9]+}", rt.getDataByID).Methods("GET", "OPTIONS")
	return
}

func New(_svc service) *Router {
	return &Router{
		svc: _svc,
	}
}
