package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PraveenBSD/content-store/models"
	"github.com/satori/go.uuid"
)

func Info(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	w.Write([]byte(fmt.Sprintf("HI %s, WELOME TO CONTENT-MANAGER API", v.Get("userName"))))

}

// ContentAccess updates content access to specific users
func ContentAccess(w http.ResponseWriter, r *http.Request) {
	log.Println("content-access endpoint hit")
	var userAccess models.UserAccess
	var errorMessage models.ErrorMessage
	err := json.NewDecoder(r.Body).Decode(&userAccess)
	if err != nil {
		log.Println(err)
		errorMessage.Code = "bad_request"
		errorMessage.Message = "userIds or contentId field missing"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	err = models.ContentAccess(&userAccess)
	if err != nil {
		log.Println(err)
		errorMessage.Code = "db_error"
		errorMessage.Message = "error updating the DB"
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	json.NewEncoder(w).Encode(userAccess)
}

// UploadContent gets the contents and stores it in contents folder
func UploadContent(w http.ResponseWriter, r *http.Request) {

	var errorMessage models.ErrorMessage
	var contentDetails models.ContentDetail

	log.Println("upload content endpoint hit")
	r.ParseMultipartForm(10 << 20)

	contentDetails.UserID = r.FormValue("userId")
	contentDetails.ID = r.FormValue("contentId")
	if contentDetails.UserID == "" {
		errorMessage.Code = "bad_request"
		errorMessage.Message = "userId field missing"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	if contentDetails.ID == "" {
		errorMessage.Code = "bad_request"
		errorMessage.Message = "contentId field missing"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		errorMessage.Code = "upload_failed"
		errorMessage.Message = fmt.Sprint("Error retrieving a File:", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	defer file.Close()
	contentDetails.Name = fmt.Sprintf("%+v", handler.Filename)
	contentDetails.Size = fmt.Sprintf("%+v", handler.Size)
	uID, err := uuid.NewV4()
	contentDetails.UploadID = fmt.Sprint(uID)

	tempFile, err := ioutil.TempFile("contents", contentDetails.UploadID+contentDetails.Name)
	if err != nil {
		log.Println(err)
		errorMessage.Code = "upload_failed"
		errorMessage.Message = fmt.Sprint("Error creating a File:", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		errorMessage.Code = "upload_failed"
		errorMessage.Message = fmt.Sprint("Error reading the File:", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	tempFile.Write(fileBytes)
	err = models.UploadContent(&contentDetails)
	if err != nil {
		log.Println(err)
		errorMessage.Code = "upload_failed"
		errorMessage.Message = fmt.Sprint("Error writing the DB:", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	json.NewEncoder(w).Encode(contentDetails)
}
