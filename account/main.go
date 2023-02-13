package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	hl "github.com/pienaahj/memrizr/account/handler"
	rt "github.com/pienaahj/memrizr/account/router"
)

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

func main() {
	port := ":8080"
	// Create a new custom router
	router := &rt.Router{}
	h := hl.NewHandler(&hl.Config{
		R: router,
	})

	log.Println("Starting server...")
	// handle all the defined routes
	h.Routes()

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
			log.Fatal("Failed to initialize server: \n", err)
		}
	}()
	/* default servmux handler
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to initialize server: \n", err)
		}
		}()
	*/

	log.Printf("Listening on port %s\n", port)

	// Wait for the kill signal
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until signal is passed to the quit channel

	<-quit

	// Shutting down server
	log.Println("Shutting down server...")
}
