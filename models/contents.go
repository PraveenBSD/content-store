package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"

	"github.com/PraveenBSD/content-store/utils"
)

// UploadContent inserts content upload details into DB
func UploadContent(content *ContentDetail) error {

	db, err := utils.Connect()
	if err != nil {
		log.Println("error connecting DB: ", err)
		return err
	}
	stmt, err := db.Prepare("INSERT INTO uploads(upload_id,content_id,content_name,content_size,user_id) VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		log.Println("error preparing DB: ", err)
		return err
	}
	_, err = stmt.Exec(content.UploadID, content.ID, content.Name, content.Size, content.UserID)
	defer db.Close()
	fmt.Println("error executing stmt: ", err)
	return err
}

// GetContentAccess - returns content-access details
func GetContentAccess(contentID string) ([]*UserAccess, error) {
	db, err := utils.Connect()
	if err != nil {
		return nil, err
	}
	var contentAccesses []*UserAccess
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s where content_id='%s'", "content_access", contentID))
	if err != nil {
		if err == sql.ErrNoRows {
			return contentAccesses, nil
		}
		return nil, err
	}
	for rows.Next() {
		var userAccess UserAccess
		var ID int
		//var arrStr string
		err = rows.Scan(&ID, &userAccess.ContentID, pq.Array(&userAccess.UserIds))
		if err != nil {
			return nil, err
		}
		contentAccesses = append(contentAccesses, &userAccess)
	}
	return contentAccesses, nil
}

// ContentAccess updates content access details into DB
func ContentAccess(userAccess *UserAccess) error {
	db, err := utils.Connect()
	if err != nil {
		log.Println("error connecting DB: ", err)
		return err
	}
	contentAccesses, err := GetContentAccess(userAccess.ContentID)
	if err != nil {
		return err
	}
	if len(contentAccesses) > 0 {
		stmt, err := db.Prepare("UPDATE content_access set user_ids=$1 where content_id=$2")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(pq.Array(userAccess.UserIds), userAccess.ContentID)
		defer db.Close()
		log.Println(err)
		return err
	} else {
		stmt, err := db.Prepare("INSERT INTO content_access(content_id, user_ids) VALUES($1,$2)")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(userAccess.ContentID, pq.Array(userAccess.UserIds))
		defer db.Close()
		log.Println(err)
		return err
	}
}
