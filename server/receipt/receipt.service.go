package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/koosie0507/pluralsight-go-webservices/server/middleware"
	"github.com/koosie0507/pluralsight-go-webservices/server/upload"
)

const receiptsPath = "receipts"

//SetupRoutes is a utility function for setting up the products API
func SetupRoutes(apiBasePath string) {
	receiptsHandler := http.HandlerFunc(handleReceipts)
	http.Handle(
		fmt.Sprintf("%s/%s", apiBasePath, receiptsPath),
		middleware.Log(middleware.JSON(middleware.CORS(receiptsHandler))),
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
