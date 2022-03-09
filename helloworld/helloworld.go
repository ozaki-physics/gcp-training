package helloworld

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mysecretmanager "github.com/ozaki-physics/gcp-training/mySecretManager"
)

func Main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Hello, World! をレスポンスするが Secret Manager の値が取得できているか確認するために GCP のログに書き込む
	projectId := "smart-ruler-277318"
	name := "test"
	if _, err := mysecretmanager.GetGCPSecretValue(projectId, name, 1); err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, "Hello, World!")
}
