package transport

import (
	"api/service"
	"be/pkg/errors"
	"net/http"

	pkghttp "be/pkg/http"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	r   chi.Router
	svc service.UserService
}

func NewUserController(r chi.Router, svc service.UserService) *UserController {
	return &UserController{r: r, svc: svc}
}

func (uc *UserController) RegisterRoutes() {
	uc.r.Post("/users/signup", uc.signup)
	uc.r.Post("/users/signin", uc.signin)
}

func (uc *UserController) signup(w http.ResponseWriter, r *http.Request) {
	var input SignUpInput
	if err := input.Bind(r); err != nil {
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
	if err := input.Bind(r); err != nil {
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
