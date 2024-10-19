package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	Router *chi.Mux
}

func NewServer() *Server {
	return &Server{Router: chi.NewRouter()}
}

func (s *Server) MountHandlers(productHandler *ProductHandler) {
	s.Router.Use(middleware.Logger)
	s.Router.Use(render.SetContentType(render.ContentTypeJSON))
	s.Router.Get("/", rootHandler)
	s.Router.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.GetAllProducts)
		r.Get("/{id}", productHandler.GetProductById)
		r.Post("/", productHandler.CreateProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Patch("/{id}", productHandler.PatchProduct)
		r.Patch("/{id}/store", productHandler.PatchStore)
	})
}
