package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aidarkhanov/nanoid"
)

type ExcalidrawDrawing struct {
	Name        string                 `json:"name"`
	DrawingJson map[string]interface{} `json:"drawingJson"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	w.Write([]byte("API is live"))
}

// func snippetView(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(r.PathValue("id"))
// 	if err != nil || id < 1 {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
// }

// func snippetCreate(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Form for creating a new snippet"))
// }

// func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write([]byte("Save a new snippet..."))
// }

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TODO: Add error handling for edge cases
// TODO: Add tests
//thebidffedfkjgjfififjfeijfjdidjd
// TODO: write compression function
func (app *application) getDrawingByName(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")
	name = strings.ReplaceAll(name, "/", "-")
	name = path.Clean(name)
	stream, err := app.dataSaver.NewReader(name)
	if err != nil {
		if os.IsNotExist(err) {
			app.clientError(w, r, http.StatusNotFound, err)
			return
		}
		app.serverError(w, r, err)
		return
	}
	defer stream.Close()

	bufferedStream := bufio.NewReader(stream)
	// if we need more control over buffer size we can change io.copy to
	//w.write and w.flush in a for loop or use second param of bufio.NewReader() function
	_, err = io.Copy(w, bufferedStream)
	if err != nil {
		app.serverError(w, r, errors.New("Error sending drawing to client:"+name))
		return
	}
}

// TODO: write post request
// TODO: make sure something over 10 MB can't be posted
// expecting body {name: string, user: username string, }
func (app *application) postDrawing(w http.ResponseWriter, r *http.Request) {
	id := nanoid.New()
	fmt.Println(id)
	// buf := bufio.NewReader(r.Body)

}

func (app *application) postCompressedDrawing(w http.ResponseWriter, r *http.Request) {
	var bodyReader io.ReadCloser
	if r.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(r.Body)
		if err != nil {
			app.serverError(w, r, errors.New("failed to decompress payload"))
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

func (app *application) streamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for i := 1; i <= 5; i++ {
		w.Write([]byte(fmt.Sprintf("Chunk %d\n", i))) // Write data to the buffer
		if flusher != nil {
			// flusher.Flush() // Immediately send the data to the client
		}
		time.Sleep(1 * time.Second) // Simulate delay between chunks
	}
}
