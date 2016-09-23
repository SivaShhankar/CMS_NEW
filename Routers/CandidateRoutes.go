package Routers

import (
	"net/http"

	handlers "github.com/SivaShhankar/CMS_NEW/Handlers"
)

func SetCandidateRoutes(router *http.ServeMux) *http.ServeMux {

	router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("Templates/css"))))
	router.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("Templates/images"))))
	router.Handle("/JS/", http.StripPrefix("/JS/", http.FileServer(http.Dir("Templates/JS"))))
	router.Handle("/Files/", http.StripPrefix("/Files/", http.FileServer(http.Dir("Templates/Files"))))

	router.Handle("/Index", http.HandlerFunc(handlers.Index))
	router.Handle("/Upload", http.HandlerFunc(handlers.Upload))
	router.Handle("/View", http.HandlerFunc(handlers.View))
	router.Handle("/Delete", http.HandlerFunc(handlers.Delete))
	router.Handle("/Search", http.HandlerFunc(handlers.Search))
	router.Handle("/Filter", http.HandlerFunc(handlers.Filter))
	router.Handle("/", http.HandlerFunc(handlers.Index))
	router.Handle("/EditData", http.HandlerFunc(handlers.Edit))

	return router
}
