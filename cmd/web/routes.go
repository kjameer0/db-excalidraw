package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	mux := http.NewServeMux()
	// the file server with assets comes from a specific folder
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// StripPrefix gets rid of /static from the URL so we aren't searching
	// for /static/static/path-to-asset
	// create a get route for all assets
	mux.Handle("GET /static/", http.StripPrefix("/static", neuter(fileServer)))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /stream/{$}", app.streamHandler)
	mux.HandleFunc("GET /drawing/{name}", app.getDrawingByName)
	mux.HandleFunc("POST /drawing/{$}", app.postDrawing)
	mux.HandleFunc("POST /compressed/drawing", app.postCompressedDrawing)
	return mux
}
