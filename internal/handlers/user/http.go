package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/domain"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/core/ports"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/server/auth"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/server/params"
	"golang.org/x/crypto/bcrypt"
)

type UserHttpHandler struct {
	userSvc ports.UserUsecase
}

func NewUserHandler(userSvc ports.UserUsecase, wsCont *restful.Container) *UserHttpHandler {
	httpHandler := &UserHttpHandler{
		userSvc: userSvc,
	}

	ws := new(restful.WebService)

	ws.Path("/user").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/register").To(httpHandler.RegisterUser))
	ws.Route(ws.POST("/login").To(httpHandler.LoginUser))
	ws.Route(ws.PUT("").To(httpHandler.UpdateUser).Filter(auth.AuthJWT))

	wsCont.Add(ws)

	return httpHandler
}

func (e *UserHttpHandler) UpdateUser(req *restful.Request, resp *restful.Response) {
	var a UpdateRequestData
	req.ReadEntity(&a)

	// get user ID for update query
	reqId, err := params.StringFrom(req.Request, auth.USER_ID_CTX_KEY)
	if err != nil || len(reqId) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("no id found for user"))
		return
	}

	ctx := req.Request.Context()
	dataUser := &domain.User{ID: reqId, Name: a.Name, Surname: a.Surname}

	err = e.userSvc.Update(ctx, dataUser)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error updating user"))
		return
	}

	// return updated user as data
	var retUser *UserModel = &UserModel{}
	retUser.FromDomain(dataUser)
	resp.WriteAsJson(retUser)
}

// Performs login or register
func (e *UserHttpHandler) RegisterUser(req *restful.Request, resp *restful.Response) {
	var reqData RegisterRequestData
	req.ReadEntity(&reqData)

	var user *UserModel = &UserModel{}
	user.Email = reqData.Email
	user.Name = reqData.Name
	user.Surname = reqData.Surname

	// todo: expand validation
	if len(user.Email) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("no email provided"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), 10)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error registering"))
		return
	}
	user.PasswordHash = string(hashedPassword)

	err = e.registerUser(req.Request.Context(), user)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error registering"))
		fmt.Println(err)
		return
	}

	authToken, err := auth.CreateJWT(user.Email, user.ID)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error creating jwt"))
		return
	}

	// send user + token back
	respData := RegisterResponseData{AuthToken: authToken, User: *user}

	resp.WriteAsJson(respData)
}

func (e *UserHttpHandler) registerUser(ctx context.Context, user *UserModel) error {
	err := e.userSvc.RegisterUser(ctx, user.ToDomain())
	if err != nil {
		return err
	}
	// retrieve their data from the DB to populate it (e.g. ID)
	userData, err := e.userSvc.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	user.FromDomain(userData)
	return nil
}

func (e *UserHttpHandler) LoginUser(req *restful.Request, resp *restful.Response) {
	var reqData LoginRequestData
	req.ReadEntity(&reqData)

	if len(reqData.Email) == 0 || len(reqData.Password) == 0 {
		resp.WriteError(http.StatusBadRequest, errors.New("bad login credentials"))
		return
	}

	userData, err := e.userSvc.FindByEmail(req.Request.Context(), reqData.Email)
	if err != nil {
		resp.WriteError(http.StatusForbidden, errors.New("unauthorized"))
		return
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(reqData.Password))
	if err != nil {
		resp.WriteError(http.StatusForbidden, errors.New("unauthorized"))
		return
	}

	authToken, err := auth.CreateJWT(userData.Email, userData.ID)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error creating jwt"))
		return
	}

	// send user + token back
	var user *UserModel = &UserModel{}
	user.FromDomain(userData)
	respData := RegisterResponseData{AuthToken: authToken, User: *user}

	resp.WriteAsJson(respData)

}
