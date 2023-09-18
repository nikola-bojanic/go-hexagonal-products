package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/config"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/usecases"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/database"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/handlers/category"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/handlers/order"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/handlers/product"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/handlers/user"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/log"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/repo"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/server/auth"
)

type Server struct {
	srv    *http.Server
	wsCont *restful.Container

	RequestLogger log.Logger
}

type ApiVersion int

const (
	V1 ApiVersion = iota
	V2
)

func NewServer(cfg config.ServerConfig, db *database.DB) *Server {

	// http Server
	httpSrv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),

		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// restful container, where all the services and routes will be connected
	wsCont := restful.NewContainer()

	// the server, encapsulating all services, logging, etc.
	fullSrv := &Server{
		srv:    httpSrv,
		wsCont: wsCont,

		RequestLogger: cfg.Logger,
	}

	// add logging
	wsCont.Filter(log.NCSACommonLogFormatLogger(cfg.Logger))

	// add error handling
	wsCont.DoNotRecover(false)
	wsCont.ServiceErrorHandler(fullSrv.WriteServiceErrorJson)
	wsCont.RecoverHandler(fullSrv.RecoverHandler)

	// base server paths
	baseWs := new(restful.WebService)
	baseWs.Path("/")
	baseWs.Route(baseWs.GET("/ping").Filter(auth.AuthJWT).To(ping))

	wsCont.Add(baseWs)

	// register routes
	userRep := repo.NewUserRepository(db)
	userSvc := usecases.NewUserService(userRep)

	categoryRep := repo.NewCategoryRepository(db)
	categorySvc := usecases.NewCategoryService(categoryRep)

	productRep := repo.NewProductRepository(db)
	productSvc := usecases.NewProductService(productRep)
	orderRep := repo.NewOrderRepository(db)
	orderSvc := usecases.NewOrderService(orderRep, productRep)
	product.NewProductHandler(productSvc, categorySvc, wsCont)
	category.NewCategoryHandler(categorySvc, wsCont)
	order.NewOrderHandler(orderSvc, productSvc, categorySvc, userSvc, wsCont)
	user.NewUserHandler(userSvc, wsCont)

	http.Handle("/", wsCont)

	return fullSrv
}

func (s *Server) ListenAndServe(env string, domain string) error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func ping(req *restful.Request, resp *restful.Response) {
	resp.Write([]byte("PONG"))
}

func (s *Server) WriteServiceErrorJson(err restful.ServiceError, req *restful.Request, resp *restful.Response) {
	s.RequestLogger.Errorf("Service error: ", err)
	resp.WriteHeader(500)
	resp.WriteAsJson("Internal server error")
}

func (s *Server) RecoverHandler(i interface{}, w http.ResponseWriter) {
	s.RequestLogger.Error("Server panic error: ", i)
	w.Write([]byte("Internal server error"))
}
