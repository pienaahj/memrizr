package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// ~~~~~ RouteEntry ~~~~~ //
// RouteEntry holds every route
type RouteEntry struct {
	Path        *regexp.Regexp
	Method      string
	HandlerFunc http.HandlerFunc
}

// parameter defines the type for key in params
type parameter string

// Match matches the path and method
func (ent *RouteEntry) Match(r *http.Request) map[parameter]string {
	match := ent.Path.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return nil // No match found
	}

	// Create a map to store URL parameters in
	params := make(map[parameter]string)
	groupNames := ent.Path.SubexpNames()
	for i, group := range match {
		params[parameter(groupNames[i])] = group
	}

	return params
}

// ~~~~~ Router ~~~~~ //
// Router struct
type Router struct {
	ctx    context.Context
	routes []RouteEntry
}

// Route routes on method, path and assigns the handlerFunc
func (rtr *Router) Route(method, path string, handlerFunc http.HandlerFunc) {
	// NOTE: ^ means start of string and $ means end. Without these,
	//   we'll still match if the path has content before or after
	//   the expression (/foo/bar/baz would match the "/bar" route).
	exactPath := regexp.MustCompile("^" + path + "$")

	e := RouteEntry{
		Method:      method,
		Path:        exactPath,
		HandlerFunc: handlerFunc,
	}
	rtr.routes = append(rtr.routes, e)
	fmt.Println("routing path: ", e.Path)
}

// ServeHTTP implements the http.Handler interface
func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR:", r)
			http.Error(w, "Uh oh!", http.StatusInternalServerError)
		}
	}()

	for _, e := range rtr.routes {
		params := e.Match(r)
		if params == nil {
			continue // No match found
		}

		// Create new request with params stored in context
		ctx := context.WithValue(r.Context(), parameter("params"), params)
		e.HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	http.NotFound(w, r)
}

// Getting the parameters
// URLParam extracts a parameter from the URL by name
func URLParam(r *http.Request, name parameter) string {
	ctx := r.Context()

	// ctx.Value returns an `interface{}` type, so we
	// also have to cast it to a map, which is the
	// type we'll be using to store our parameters.
	params := ctx.Value("params").(map[parameter]string)
	return params[name]
}
