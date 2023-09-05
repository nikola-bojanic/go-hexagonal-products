package order

import (
	"errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
)

type OrderHttpHandler struct {
	orderSvc ports.OrderUsecase
}

func NewOrderHandler(orderSvc ports.OrderUsecase, wsCont *restful.Container) *OrderHttpHandler {
	httpHandler := &OrderHttpHandler{
		orderSvc: orderSvc,
	}

	ws := new(restful.WebService)
	ws.Path("/order").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.POST("/").To(httpHandler.CreateOrder))
	ws.Route(ws.POST("/status").To(httpHandler.UpdateOrderStatus))
	ws.Route(ws.PUT("/status").To(httpHandler.UpdateOrderStatus))
	ws.Route(ws.DELETE("/status").To(httpHandler.UpdateOrderStatus))

	wsCont.Add(ws)

	return httpHandler
}

func (e *OrderHttpHandler) CreateOrder(req *restful.Request, res *restful.Response) {
	var reqData OrderRequest
	req.ReadEntity(&reqData)
	var order *OrderModel = &OrderModel{}
	order.Status = reqData.Status
	order.ProductItems = reqData.Products
	created, err := e.orderSvc.CreateOrder(req.Request.Context(), order.ToDomain())
	if err != nil {
		res.WriteError(http.StatusInternalServerError, errors.New("error creating order"))
		return
	}
	order.FromDomain(created)
	res.WriteAsJson(order)
}

func (e *OrderHttpHandler) UpdateOrderStatus(req *restful.Request, res *restful.Response) {
	var reqData OrderRequest
	req.ReadEntity(&reqData)
	var order *OrderModel = &OrderModel{}
	order.ID = reqData.OrderId
	order.Status = reqData.Status
	updated, err := e.orderSvc.UpdateOrderStatus(req.Request.Context(), order.ToDomain())
	if err != nil {
		res.WriteError(http.StatusInternalServerError, errors.New("error updating order status"))
		return
	}
	order.FromDomain(updated)
	res.WriteAsJson(order)
}
