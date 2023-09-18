package order

import (
	"errors"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/server/auth"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/server/params"
)

type OrderHttpHandler struct {
	orderSvc    ports.OrderUsecase
	productSvc  ports.ProductUsecase
	categorySvc ports.CategoryUsecase
	userSvc     ports.UserUsecase
}

func NewOrderHandler(orderSvc ports.OrderUsecase, productSvc ports.ProductUsecase, categorySvc ports.CategoryUsecase, userSvc ports.UserUsecase, wsCont *restful.Container) *OrderHttpHandler {
	httpHandler := &OrderHttpHandler{
		orderSvc:    orderSvc,
		productSvc:  productSvc,
		categorySvc: categorySvc,
		userSvc:     userSvc,
	}

	ws := new(restful.WebService)
	ws.Path("/order").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.POST("/").To(httpHandler.CreateOrder).Filter(auth.AuthJWT))
	ws.Route(ws.PUT("/").To(httpHandler.UpdateOrderStatus).Filter(auth.AuthJWT))
	ws.Route(ws.DELETE("/").To(httpHandler.DeleteOrder))

	wsCont.Add(ws)

	return httpHandler
}

func (e *OrderHttpHandler) CreateOrder(req *restful.Request, res *restful.Response) {
	var reqData OrderRequest
	req.ReadEntity(&reqData)
	reqId, err := params.StringFrom(req.Request, auth.USER_ID_CTX_KEY)
	if err != nil || len(reqId) == 0 {
		res.WriteError(http.StatusBadRequest, errors.New("no id found for user"))
		return
	}

	var order *OrderModel = &OrderModel{}
	order.User.ID = reqId
	order.Status = reqData.Status
	order.ProductItems = reqData.Products
	created, err := e.orderSvc.CreateOrder(req.Request.Context(), order.ToDomain())
	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}
	order.FromDomain(created)
	res.WriteAsJson(order)
}

func (e *OrderHttpHandler) UpdateOrderStatus(req *restful.Request, res *restful.Response) {
	var reqData OrderRequest
	req.ReadEntity(&reqData)
	reqId, err := params.StringFrom(req.Request, auth.USER_ID_CTX_KEY)
	if err != nil || len(reqId) == 0 {
		res.WriteError(http.StatusBadRequest, errors.New("no id found for user"))
		return
	}
	toUpdate, err := e.orderSvc.FindOrderById(req.Request.Context(), reqData.ID)
	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}
	if toUpdate.ID != reqData.ID {
		res.WriteError(http.StatusBadRequest, errors.New("user cannot edit other user's order"))
		return
	}
	var order *OrderModel = &OrderModel{}
	order.ID = reqData.ID
	order.Status = reqData.Status
	order.ProductItems = reqData.Products
	order.User.ID = reqId
	updated, err := e.orderSvc.UpdateOrderStatus(req.Request.Context(), order.ToDomain())
	if err != nil {
		res.WriteError(http.StatusInternalServerError, errors.New("error updating order"))
		return
	}
	order.FromDomain(updated)
	res.WriteAsJson(order)
}

func (e *OrderHttpHandler) DeleteOrder(req *restful.Request, res *restful.Response) {
	var reqData OrderRequest
	req.ReadEntity(&reqData)
	var order *OrderModel = &OrderModel{}
	order.ID = reqData.ID
	order.ProductItems = reqData.Products
	order.Status = reqData.Status
	err := e.orderSvc.DeleteOrder(req.Request.Context(), order.ToDomain())
	if err != nil {
		res.WriteError(http.StatusInternalServerError, errors.New("error deleting order"))
		return
	}
	res.WriteAsJson(order)
}
