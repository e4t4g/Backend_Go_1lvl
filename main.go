package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func main() {
	uploadHandler := &UploadHandler{
		UploadDir: "upload",
		HostAddr:  "localhost:8001",
	}

	//dirToServe := http.Dir(uploadHandler.UploadDir)

	fs := &http.Server{
		Addr:         ":8001",
		Handler:      nil, //http.FileServer(dirToServe)
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.Handle("/upload", uploadHandler)

	newSrv := &http.Server{
		Addr:         ":8002",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.HandleFunc("/list", uploadHandler.listGetFiles)

	go func() {
		log.Fatal(fs.ListenAndServe())
	}()

	log.Fatal(newSrv.ListenAndServe())

}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	filePath := h.UploadDir + "/" + header.Filename

	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File %s has been successfully uploaded", header.Filename)

	fileLink := h.HostAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)

}

func (h *UploadHandler) listGetFiles(w http.ResponseWriter, r *http.Request) {
	dirFiles, err := ioutil.ReadDir(h.UploadDir)
	if err != nil {
		log.Fatal(err)
	}

	ext := r.FormValue("ext")

	if len(ext) == 0 {
		for _, files := range dirFiles {
			fmt.Fprintln(w,
				files.Name(),
				filepath.Ext(filepath.Ext(files.Name())),
				files.Size())
		}
	} else {
		for _, files := range dirFiles {
			if filepath.Ext(filepath.Ext(files.Name())) == ext {
				fmt.Fprintln(w,
					files.Name(),
					filepath.Ext(filepath.Ext(files.Name())),
					files.Size())
			}
		}
	}

}
