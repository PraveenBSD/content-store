package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PraveenBSD/content-store/models"
	"github.com/PraveenBSD/content-store/utils"
)

const (
	// ServerPrivateKeyPath - server private key path
	ServerPrivateKeyPath = "/Users/praveen/go/src/github.com/PraveenBSD/content-store/keys/id_rsa"
	// UserCertificatePath - user certificate folder path
	UserCertificatePath = "/Users/praveen/go/src/github.com/PraveenBSD/content-store/keys/user_keys/"
	// ContentsPath - contents folder path
	ContentsPath = "/Users/praveen/go/src/github.com/PraveenBSD/content-store/contents/"
)

func isContentAccessableByUser(contentID string, userID string) (bool, error) {
	contentAccesses, err := models.GetContentAccess(contentID)
	if err != nil {
		log.Println(err)
		return false, err
	}
	for _, contentAccess := range contentAccesses {
		for _, contentUserID := range contentAccess.UserIds {
			if userID == contentUserID {
				return true, nil
			}
		}
	}
	return false, nil

}

// DownloadContent - API to download content
func DownloadContent(w http.ResponseWriter, r *http.Request) {
	log.Println("download content endpoint hit")
	var downloadContent models.DownloadContent
	var errorMessage models.ErrorMessage
	err := json.NewDecoder(r.Body).Decode(&downloadContent)
	if err != nil {
		errorMessage.Code = "bad_request"
		errorMessage.Message = "contentId or userId field missing"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	isAccessible, err := isContentAccessableByUser(downloadContent.ContentID, downloadContent.UserID)
	if err != nil {
		errorMessage.Code = "unexpected_error"
		errorMessage.Message = fmt.Sprint("Error: ", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	if isAccessible {
		encrFilePath := ContentsPath + downloadContent.ContentID
		contentServerEncr, err := ioutil.ReadFile(encrFilePath)
		if err != nil {
			errorMessage.Code = "unexpected_error"
			errorMessage.Message = fmt.Sprint("Error: ", err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(errorMessage)
			return
		}
		content, err := utils.DecryptFile(ServerPrivateKeyPath, contentServerEncr)
		if err != nil {
			log.Println(err)
			errorMessage.Code = "unexpected_error"
			errorMessage.Message = fmt.Sprint("Error in decrypting: ", err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(errorMessage)
			return
		}
		contentUserEncr, err := utils.EncryptFile(UserCertificatePath+downloadContent.UserID+".pub.pem", content)
		if err != nil {
			log.Println(err)
			errorMessage.Code = "unexpected_error"
			errorMessage.Message = fmt.Sprint("Error in encrypting: ", err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(errorMessage)
			return
		}
		downloadContent.Content = base64.StdEncoding.EncodeToString(contentUserEncr)
		json.NewEncoder(w).Encode(downloadContent)
	} else {
		log.Println("unauthorized user error")
		errorMessage.Code = "unauthorized_error"
		errorMessage.Message = "user doesnt have access to the content"
		w.WriteHeader(503)
		json.NewEncoder(w).Encode(errorMessage)
	}

}
