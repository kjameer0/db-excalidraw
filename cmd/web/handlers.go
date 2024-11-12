package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type ExcalidrawDrawing struct {
	Name        string                 `json:"name"`
	DrawingJson map[string]interface{} `json:"drawingJson"`
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	// base template must always come first in slice
	TestCompression()
	// files := []string{"./ui/html/base.tmpl.html", "./ui/html/pages/home.tmpl", "./ui/html/partials/nav.tmpl"}
	// templateSet, err := template.ParseFiles(files...)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	// // nil in Execute means there is no custom data to add to the template
	// err = templateSet.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
	w.Write([]byte("compression"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Form for creating a new snippet"))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getDrawingByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	name = strings.ReplaceAll(name, "/", "-")
	name = path.Clean(name)
	pathName := "./drawings/" + name + ".json"
	_, err := os.Stat(pathName)
	if os.IsNotExist(err) {
		slog.Info("file does not exist GET /drawing/{name} " + name)
		http.NotFound(w, r)
		return
	}
	drawing, err := os.ReadFile(pathName)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Write([]byte(drawing))
}

func postDrawing(w http.ResponseWriter, r *http.Request) {
	// respond with error if file exists
	drawing := ExcalidrawDrawing{}
	err := json.NewDecoder(r.Body).Decode(&drawing)
	drawing.Name = strings.ReplaceAll(drawing.Name, "/", "-")
	// if file json cannot be parsed
	if err != nil {
		slog.Error(err.Error(), "reason", "failed to parse json")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !validateExcalidrawDrawing(&drawing) {
		slog.Error("failed to validate excalidraw drawing")
		http.Error(w, "invalid excalidraw drawing", http.StatusInternalServerError)
		return
	}
	drawingPathName := "./drawings/" + drawing.Name + ".json"
	var file *os.File = nil
	// check for file existence
	_, err = os.Stat(drawingPathName)
	// file doesn't exist, create
	if err != nil {
		file, err = os.Create(drawingPathName)
		if err != nil {
			fmt.Println(err)
			slog.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// file exists already return error
		slog.Info("file exists", slog.String("fileName", drawingPathName))
		http.Error(w, "name for drawing already exists", http.StatusConflict)
		return
	}
	defer file.Close()

	marshaledDrawing, err := json.Marshal(&drawing)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = file.Write(marshaledDrawing)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("success"))
}

func postCompressedDrawing(w http.ResponseWriter, r *http.Request) {
	var bodyReader io.ReadCloser
	if r.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(r.Body)
		if err != nil {
			slog.Error(err.Error(), "reason", "Failed to decompress payload")
			http.Error(w, "Failed to decompress payload", http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()
		bodyReader = gzipReader
	} else {
		bodyReader = r.Body
	}
	defer bodyReader.Close()

	contentLength := r.ContentLength
	if contentLength == 0 {
		slog.Error("Request body is empty", "reason", "Empty payload")
		http.Error(w, "Empty payload", http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		slog.Error(err.Error(), "reason", "Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	fmt.Println((string(bodyBytes)[:400]))

	// respond with error if file exists
	drawing := ExcalidrawDrawing{}
	err = json.Unmarshal(bodyBytes, &drawing)
	drawing.Name = strings.ReplaceAll(drawing.Name, "/", "-")
	if err != nil {
		slog.Error(err.Error(), "reason", "failed to parse json")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !validateExcalidrawDrawing(&drawing) {
		slog.Error("failed to validate excalidraw drawing")
		http.Error(w, "invalid excalidraw drawing", http.StatusInternalServerError)
		return
	}
	drawingPathName := "./drawings/" + drawing.Name + ".json"
	var file *os.File = nil
	// check for file existence
	_, err = os.Stat(drawingPathName)
	// file doesn't exist, create
	if err != nil {
		file, err = os.Create(drawingPathName)
		if err != nil {
			fmt.Println(err)
			slog.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// file exists already return error
		slog.Info("file exists", slog.String("fileName", drawingPathName))
		http.Error(w, "name for drawing already exists", http.StatusConflict)
		return
	}
	defer file.Close()

	unMarshaledDrawing, err := json.Marshal(&drawing)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = file.Write(unMarshaledDrawing)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//make more specific response
	w.Write([]byte("success"))
}

func validateExcalidrawDrawing(drawingInfo *ExcalidrawDrawing) bool {
	fmt.Printf("Drawing %s\n****\n", drawingInfo.DrawingJson)
	drawing := drawingInfo.DrawingJson
	drawingType, ok := drawing["type"]
	if !ok {
		return false
	}

	drawingTypeString, ok := drawingType.(string)
	if !ok {
		slog.Error("drawing.type is not a string")
		return false
	}

	if !(strings.Contains(drawingTypeString, "excalidraw")) {
		return false
	}

	elements, ok := drawing["elements"]
	if !ok {
		return false
	}

	_, ok = elements.([]interface{})
	return ok
}
