package web

import (
	"librarymanagement/web/handlers"
	"librarymanagement/web/middlewire"
	"net/http"
)

func InitRoutes(mux *http.ServeMux, manager *middlewire.Manager) {
	mux.Handle(
		"POST /reader/register",
		manager.With(
			http.HandlerFunc(handlers.RegisterReader),
		),
	)
	mux.Handle(
		"POST /reader/login",
		manager.With(
			http.HandlerFunc(handlers.LoginReader),
		),
	)

	mux.Handle(
		"GET /admin/unapproveduser",
		manager.With(
			http.HandlerFunc(handlers.GetUnapprovedUser), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"GET /admin/fetchuser",
		manager.With(
			http.HandlerFunc(handlers.FetchUser), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"GET /admin/fetchadmin",
		manager.With(
			http.HandlerFunc(handlers.FetchAdmin), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"DELETE /admin/deleteadmin",
		manager.With(
			http.HandlerFunc(handlers.DeleteAdmin), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"POST /admin/acceptapproval",
		manager.With(
			http.HandlerFunc(handlers.ApprovedUser), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"POST /admin/login",
		manager.With(
			http.HandlerFunc(handlers.LoginAdmin),
		),
	)

	mux.Handle(
		"POST /admin/addadmin",
		manager.With(
			http.HandlerFunc(handlers.AddAdmin), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"POST /admin/addbook",
		manager.With(
			http.HandlerFunc(handlers.AddBook), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"PUT /admin/updatebook/{id}",
		manager.With(
			http.HandlerFunc(handlers.UpdateBook), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"DELETE /admin/deletebook",
		manager.With(
			http.HandlerFunc(handlers.DeleteBook), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"DELETE /admin/delete-request",
		manager.With(
			http.HandlerFunc(handlers.DeleteBorrowRequest), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"GET /reader/searchbook",
		manager.With(
			http.HandlerFunc(handlers.SearchBook), middlewire.AuthenticateUser,
		),
	)

	mux.Handle(
		"GET /admin/searchbook",
		manager.With(
			http.HandlerFunc(handlers.SearchBook), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"POST /reader/bookrequest",
		manager.With(
			http.HandlerFunc(handlers.BorrowRequestBook), middlewire.AuthenticateUser,
		),
	)

	mux.Handle(
		"GET /admin/borrowedbook",
		manager.With(
			http.HandlerFunc(handlers.FetchBorrowStatus), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"PATCH /admin/approvedbookrequest",
		manager.With(
			http.HandlerFunc(handlers.ApprovedBookRequest), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"PATCH /admin/rejectborrowreq",
		manager.With(
			http.HandlerFunc(handlers.RejectRequest), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"GET /reader/history",
		manager.With(
			http.HandlerFunc(handlers.UserHistory), middlewire.AuthenticateUser,
		),
	)

	mux.Handle(
		"GET /admin/user-history",
		manager.With(
			http.HandlerFunc(handlers.UserHistory), middlewire.AuthenticateAdmin,
		),
	)

	mux.Handle(
		"PATCH /reader/updatereadprogress",
		manager.With(
			http.HandlerFunc(handlers.UserReadProgressUpdate), middlewire.AuthenticateUser,
		),
	)

	mux.Handle(
		"PATCH /reader/returnbook",
		manager.With(
			http.HandlerFunc(handlers.ReturnBook), middlewire.AuthenticateUser,
		),
	)

	mux.Handle(
		"GET /reader/userinfo",
		manager.With(
			http.HandlerFunc(handlers.UserInfo), middlewire.AuthenticateUser,
		),
	)

	mux.Handle(
		"GET /admin/user-info",
		manager.With(
			http.HandlerFunc(handlers.UserInfo), middlewire.AuthenticateAdmin,
		),
	)
}
