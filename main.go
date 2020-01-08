package main

import (
	"fmt"
	"net/http"
	"uploadSystem/handler"
)

func main() {
	http.HandleFunc("/fileupload", handler.UploadHandle)
	http.HandleFunc("/upload/success", handler.UploadSuccess)
	http.HandleFunc("/upload/get", handler.GetFileHandler)
	http.HandleFunc("/download", handler.Download)
	http.HandleFunc("/rename", handler.RenameFile)
	http.HandleFunc("/delete", handler.DeleteFile)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Printf("%v", err)
	}
}
