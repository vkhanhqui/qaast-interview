package transport

import (
	"api/service"
	"be/pkg/errors"
	"encoding/json"
	"net/http"

	pkghttp "be/pkg/http"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	r      chi.Router
	svc    service.UserService
	jwtKey string
}

func NewUserController(r chi.Router, svc service.UserService, jwtKey string) *UserController {
	return &UserController{r: r, svc: svc, jwtKey: jwtKey}
}

func (uc *UserController) RegisterRoutes() {
	uc.r.Post("/users/signup", uc.signup)
	uc.r.Post("/users/signin", uc.signin)

	uc.r.Group(func(r chi.Router) {
		r.Use(pkghttp.AuthMiddleware(uc.jwtKey))
		r.Put("/users", uc.updateUser)
	})
}

func (uc *UserController) signup(w http.ResponseWriter, r *http.Request) {
	var input SignUpInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkghttp.JSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	id, err := uc.svc.SignUp(r.Context(), input.Email, input.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.IsInvalid(err) {
			status = http.StatusBadRequest
		}
		pkghttp.JSON(w, status, ErrorResponse{Error: err.Error()})
		return
	}

	pkghttp.JSON(w, http.StatusCreated, SignUpResponse{UserID: id})
}

func (uc *UserController) signin(w http.ResponseWriter, r *http.Request) {
	var input SignInInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkghttp.JSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	token, _, err := uc.svc.SignIn(r.Context(), input.Email, input.Password)
	if err != nil {
		pkghttp.JSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	pkghttp.JSON(w, http.StatusOK, SignInResponse{Token: token})
}

func (uc *UserController) updateUser(w http.ResponseWriter, r *http.Request) {
	uid := pkghttp.GetUserID(w, r)
	if uid == "" {
		return
	}

	var input UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkghttp.JSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	u, err := uc.svc.UpdateUser(r.Context(), uid, input.Email, input.Name)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.IsInvalid(err) {
			status = http.StatusBadRequest
		}
		pkghttp.JSON(w, status, ErrorResponse{Error: err.Error()})
		return
	}

	res := UpdateUserResponse{}
	res.Bind(u)
	pkghttp.JSON(w, http.StatusOK, res)
}
