package server

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Solbreen/final-project/pkg/api"
)

func Run() {
	port := flag.String("port", "7540", "port to serve on")
	flag.Parse()

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		*port = envPort
	}

	//mux := http.NewServeMux()

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

	fmt.Println("Запускаем сервер")
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
}
