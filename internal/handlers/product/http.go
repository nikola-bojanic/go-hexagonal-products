package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/handlers/category"
)

type ProductHttpHandler struct {
	productSvc  ports.ProductUsecase
	categorySvc ports.CategoryUsecase
}

func NewProductHandler(productSvc ports.ProductUsecase, categorySvc ports.CategoryUsecase, wsCont *restful.Container) *ProductHttpHandler {
	httpHandler := &ProductHttpHandler{
		productSvc:  productSvc,
		categorySvc: categorySvc,
	}

	ws := new(restful.WebService)

	ws.Path("/product").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(httpHandler.GetProducts))
	ws.Route(ws.GET("/{id}").To(httpHandler.GetProduct))
	ws.Route(ws.POST("").To(httpHandler.CreateProduct))
	ws.Route(ws.DELETE("/{id}").To(httpHandler.DeleteProduct))
	ws.Route(ws.PUT("/{id}").To(httpHandler.UpdateProduct))

	wsCont.Add(ws)

	return httpHandler
}

func (e *ProductHttpHandler) GetProducts(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()
	products, err := e.productSvc.GetAllProducts(ctx)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error retrieving products"))
		return
	}
	var retProducts []ProductModel
	var retProduct *ProductModel = &ProductModel{Category: &category.CategoryModel{}}

	for _, product := range *products {
		retProduct.FromDomain(&product)
		retProducts = append(retProducts, *retProduct)
	}
	resp.WriteAsJson(retProducts)
}

func (e *ProductHttpHandler) GetProduct(req *restful.Request, resp *restful.Response) {
	id, err := getId(req, resp)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, errors.New("invalid product id"))
		return
	}
	product, err := e.productSvc.FindProductById(req.Request.Context(), id)
	if err != nil {
		resp.WriteError(http.StatusNotFound, errors.New("product doesn't exist"))
		return
	}
	var retProduct *ProductModel = &ProductModel{Category: &category.CategoryModel{}}
	retProduct.FromDomain(product)
	resp.WriteAsJson(retProduct)
}

func (e *ProductHttpHandler) CreateProduct(req *restful.Request, resp *restful.Response) {
	var reqData ProductRequest
	req.ReadEntity(&reqData)

	var product *ProductModel = &ProductModel{}
	product.Name = reqData.Name
	product.Description = reqData.Description
	product.ShortDescription = reqData.ShortDescription
	product.Quantity = reqData.Quantity
	product.Price = reqData.Price
	userCategory, err := e.categorySvc.FindCategoryById(req.Request.Context(), int64(reqData.Category.Id))
	if err != nil {
		resp.WriteError(http.StatusNotFound, errors.New("category doesn't exist"))
		return
	}
	var category *category.CategoryModel = &category.CategoryModel{}
	category.FromDomain(userCategory)
	product.Category = category
	id, err := e.productSvc.CreateProduct(req.Request.Context(), product.ToDomain())
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error creating product"))
		return
	}
	resp.WriteAsJson(Response{ID: id, Message: "product created"})
}

func (e *ProductHttpHandler) DeleteProduct(req *restful.Request, resp *restful.Response) {
	id, err := getId(req, resp)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, errors.New("invalid product id"))
		return
	}
	deleted, err := e.productSvc.DeleteProduct(req.Request.Context(), id)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("an error occured"))
		return
	}
	if deleted == 0 {
		resp.WriteError(http.StatusNotFound, errors.New("product doesn't exist"))
		return
	}
	resp.WriteAsJson(Response{ID: deleted, Message: "product deleted"})
}

func (e *ProductHttpHandler) UpdateProduct(req *restful.Request, resp *restful.Response) {
	id, err := getId(req, resp)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, errors.New("invalid product id"))
		return
	}
	var productReq ProductRequest
	req.ReadEntity(&productReq)
	userCategory, err := e.categorySvc.FindCategoryById(req.Request.Context(), int64(productReq.Category.Id))
	if err != nil {
		resp.WriteError(http.StatusNotFound, errors.New("category doesn't exist"))
		return
	}
	dataProduct := &domain.Product{Name: productReq.Name, ShortDescription: productReq.ShortDescription, Description: productReq.Description,
		Quantity: productReq.Quantity, Price: productReq.Price, Category: userCategory}
	updated, err := e.productSvc.UpdateProduct(req.Request.Context(), dataProduct, id)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("an error occured"))
		return
	}
	if updated == 0 {
		resp.WriteError(http.StatusNotFound, errors.New("product doesn't exist"))
		return
	}
	resp.WriteAsJson(Response{ID: updated, Message: "product updated"})
}

func getId(req *restful.Request, resp *restful.Response) (int64, error) {
	idS := req.PathParameter("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}
