package transport

import (
	"api/service"
	"net/http"

	"be/pkg/errors"
	pkghttp "be/pkg/http"

	"github.com/go-chi/chi/v5"
)

type AdminController struct {
	r        chi.Router
	userSvc  service.UserService
	adminSvc service.AdminService
	jwtKey   string
}

func NewAdminController(r chi.Router, userSvc service.UserService, adminSvc service.AdminService, jwtKey string) *AdminController {
	return &AdminController{r: r, userSvc: userSvc, adminSvc: adminSvc, jwtKey: jwtKey}
}

func (uc *AdminController) RegisterRoutes() {
	uc.r.Group(func(r chi.Router) {
		r.Use(pkghttp.AuthMiddleware(uc.jwtKey))
		r.Get("/admin/users", uc.listUsers)
		r.Put("/admin/users", uc.updateUser)

		r.Get("/admin/userlogs", uc.listUserLogs)
	})
}

func (uc *AdminController) listUsers(w http.ResponseWriter, r *http.Request) {
	input := AdminListUsersInput{}
	input.Bind(r.URL.Query())

	users, err := uc.adminSvc.ListUsers(r.Context(), input.Limit, input.Cursor)
	if err != nil {
		pkghttp.JSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	res := AdminListUsersResponse{}
	res.Bind(users)
	pkghttp.JSON(w, http.StatusOK, res)
}

func (uc *AdminController) listUserLogs(w http.ResponseWriter, r *http.Request) {
	input := AdminListUserLogsInput{}
	input.Bind(r.URL.Query())

	userLogs, cursor, err := uc.adminSvc.ListUserLogs(r.Context(), input.Limit, input.Cursor)
	if err != nil {
		pkghttp.JSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	res := AdminListUserLogsResponse{}
	res.Bind(userLogs, cursor)
	pkghttp.JSON(w, http.StatusOK, res)
}

func (uc *AdminController) updateUser(w http.ResponseWriter, r *http.Request) {
	var input AdminUpdateUserInput
	if err := input.Bind(r); err != nil {
		pkghttp.JSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	u, err := uc.userSvc.UpdateUser(r.Context(), input.ID, input.Email, input.Name)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.IsInvalid(err) {
			status = http.StatusBadRequest
		}
		if errors.IsNotFound(err) {
			status = http.StatusNotFound
		}
		pkghttp.JSON(w, status, ErrorResponse{Error: err.Error()})
		return
	}

	res := AdminUpdateUserResponse{}
	res.Bind(u)
	pkghttp.JSON(w, http.StatusOK, res)
}
