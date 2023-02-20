package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	// rt "github.com/pienaahj/memrizr/account/router"
)

func main() {
	port := ":8080"

	// initialize data sources
	ds, err := initDS()

	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}
	// Create a new custom router
	// router := &rt.Router{}
	router, err := inject(ds)
	if err != nil {
		log.Fatalf("Failure to inject data sources: %v\n", err)
	}
	// h := hl.NewHandler(&hl.Config{
	// 	R: router,
	// })
	// handler.NewHandler(&handler.Config{
	// 	R: router,
	// })
	log.Println("Starting server...")
	// handle all the defined routes
	// h.Routes()

	// Create a custom http server with a custom handler(router)
	s := &http.Server{
		Addr:         port,
		Handler:      router, // router
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Gracefull shutdown
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()
	/* default servmux handler
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: \n", err)
		}
		}()
	*/

	log.Printf("Listening on port %s\n", port)

	// Wait for the kill signal
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until signal is passed to the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources
	if err := ds.close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	// Shutting down server
	log.Println("Shutting down server...")
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

/*				own router use

				r := &Router{}

				r.Route(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("The Best Router!"))
				})

				r.Route(http.MethodGet, `/hello/(?P<Message>\w+)`, func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("Hello " + URLParam(r, "Message")))
				})

				r.Route(http.MethodGet, "/panic", func(w http.ResponseWriter, r *http.Request) {
					panic("something bad happened!")
				})

				http.ListenAndServe(":8000", r)

*/
