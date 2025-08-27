package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// http.HandleFunc("/download-review", handlers.DownloadReviewHandler)

	port := 8080
	fmt.Printf("Microservice running at http://localhost:%d\n", port)
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
