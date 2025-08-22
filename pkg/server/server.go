package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Solbreen/final-project/pkg/api"
)

func Run() {

	port := flag.String("port", "7540", "port to serve on")
	flag.Parse()

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		*port = envPort
	}

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: nil,
	}

	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})

	api.Init()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Запускаем сервер. Порт:", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ошибка сервера: %v\n", err)
			os.Exit(1)
		}
	}()

	<-stop
	fmt.Println("\nСтоп сервер")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Ошибка при остановке сервера: %v\n", err)
		os.Exit(1)
	}
}
