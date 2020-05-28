package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/koosie0507/pluralsight-go-webservices/server/middleware"
	"github.com/koosie0507/pluralsight-go-webservices/server/upload"
)

const receiptsPath = "receipts"

//SetupRoutes is a utility function for setting up the products API
func SetupRoutes(apiBasePath string) {
	receiptsHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)
	http.Handle(
		fmt.Sprintf("%s/%s", apiBasePath, receiptsPath),
		middleware.Log(middleware.CORS(receiptsHandler)),
	)
	http.Handle(
		fmt.Sprintf("%s/%s/", apiBasePath, receiptsPath),
		middleware.Log(middleware.CORS(downloadHandler)),
	)
}

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		files, err := upload.GetUploads()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		j, err := json.Marshal(files)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		r.ParseMultipartForm(5 << 20) // 5 MB
		file, handler, err := r.FormFile("receipt")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()
		if _, err := os.Stat(upload.UploadsDir); os.IsNotExist(err) {
			os.Mkdir(upload.UploadsDir, 0755)
		}
		f, err := os.OpenFile(filepath.Join(upload.UploadsDir, handler.Filename), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	urlPathSegments := strings.Split(r.URL.Path, receiptsPath+"/")
	tail := urlPathSegments[1:]
	if len(tail) > 1 {
		w.WriteHeader(http.StatusBadRequest)
	}
	fileName := tail[0]
	file, err := os.Open(filepath.Join(upload.UploadsDir, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	contentType := http.DetectContentType(fileHeader)
	stats, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fileSize := strconv.FormatInt(stats.Size(), 10)
	w.Header().Set("Content-Disposition", "attachement; filename="+fileName)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fileSize)
	file.Seek(0, 0)
	io.Copy(w, file)
}
