package category

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
)

type CategoryHttpHandler struct {
	categorySvc ports.CategoryUsecase
}

func NewCategoryHandler(categorySvc ports.CategoryUsecase, wsCont *restful.Container) *CategoryHttpHandler {
	httpHandler := &CategoryHttpHandler{
		categorySvc: categorySvc,
	}

	ws := new(restful.WebService)

	ws.Path("/category").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(httpHandler.GetCategories))
	ws.Route(ws.GET("/{id}").To(httpHandler.GetCategory))
	ws.Route(ws.POST("").To(httpHandler.CreateCategory))
	ws.Route(ws.DELETE("/{id}").To(httpHandler.DeleteCategory))
	ws.Route(ws.PUT("/{id}").To(httpHandler.UpdateCategory))

	wsCont.Add(ws)

	return httpHandler
}

func (e *CategoryHttpHandler) GetCategories(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()
	categories, err := e.categorySvc.GetAllCategories(ctx)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error retrieving categories"))
		return
	}
	var retCategories []CategoryModel
	var retCategory *CategoryModel = &CategoryModel{}

	for _, category := range *categories {
		retCategory.FromDomain(&category)
		retCategories = append(retCategories, *retCategory)
	}
	resp.WriteAsJson(retCategories)
}

func (e *CategoryHttpHandler) GetCategory(req *restful.Request, resp *restful.Response) {
	id, err := getId(req, resp)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, errors.New("invalid category id"))
		return
	}
	category, err := e.categorySvc.FindCategoryById(req.Request.Context(), id)
	if err != nil {
		resp.WriteError(http.StatusNotFound, errors.New("category doesn't exist"))
		return
	}
	var retCategory *CategoryModel = &CategoryModel{}
	retCategory.FromDomain(category)
	resp.WriteAsJson(retCategory)
}

func (e *CategoryHttpHandler) CreateCategory(req *restful.Request, resp *restful.Response) {
	var reqData CategoryRequest
	req.ReadEntity(&reqData)

	var category *CategoryModel = &CategoryModel{}
	category.Name = reqData.Name

	if len(category.Name) < 1 {
		resp.WriteError(http.StatusBadRequest, errors.New("name not provided"))
		return
	}
	categoryId, err := e.categorySvc.CreateCategory(req.Request.Context(), category.ToDomain())

	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error creating category"))
		return
	}
	resp.WriteAsJson(Response{ID: categoryId, Name: category.Name})
}

func (e *CategoryHttpHandler) DeleteCategory(req *restful.Request, resp *restful.Response) {
	id, err := getId(req, resp)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, errors.New("invalid category id"))
		return
	}
	rows, err := e.categorySvc.DeleteCategory(req.Request.Context(), id)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("an error occured"))
		return
	}
	if rows == 0 {
		resp.WriteError(http.StatusNotFound, errors.New("category doesn't exist"))
		return
	}
	resp.WriteAsJson(Response{ID: rows, Name: "category deleted"})
}

func (e *CategoryHttpHandler) UpdateCategory(req *restful.Request, resp *restful.Response) {
	id, err := getId(req, resp)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, errors.New("invalid category id"))
		return
	}
	var categoryReq CategoryRequest
	req.ReadEntity(&categoryReq)
	dataCategory := &domain.Category{Name: categoryReq.Name}
	updated, err := e.categorySvc.UpdateCategory(req.Request.Context(), dataCategory, id)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("an error occured"))
		return
	}
	if updated == 0 {
		resp.WriteError(http.StatusNotFound, errors.New("category doesn't exist"))
		return
	}
	resp.WriteAsJson(Response{ID: updated, Name: dataCategory.Name})

}

func getId(req *restful.Request, resp *restful.Response) (int64, error) {
	idS := req.PathParameter("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}
