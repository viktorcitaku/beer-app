package router

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/viktorcitaku/beer-app/internal/controller/v1"
)

func Router(controller *v1.Controller, path string) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Get("/api/v1/hello-world", controller.HelloWorld)
	mux.Post("/api/v1/user-profiles", controller.SaveUserProfiles)
	mux.Get("/api/v1/beers", controller.GetBeers)
	mux.Post("/api/v1/beers", controller.SaveBeers)
	mux.Get("/api/v1/user-preferences", controller.GetUserPreferences)
	mux.Post("/api/v1/user-preferences", controller.SaveUserPreferences)

	FileServer(mux, "/", http.Dir(path))

	return mux
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		routePattern := ctx.RoutePattern()
		pathPrefix := strings.TrimSuffix(routePattern, "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
