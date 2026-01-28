package main

import (
	"net/http"
	"shareapp/internal/data"
	"shareapp/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func (app *application) handleMediaPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		r.Body = http.MaxBytesReader(w, r.Body, 500<<20)

		file, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		defer file.Close()

		fileName := fileHeader.Filename
		fileSize := fileHeader.Size
		contentType := fileHeader.Header.Get("Content-Type")

		userID, ok := r.Context().Value("userID").(uuid.UUID)
		if !ok {
			app.serverErrorResponse(w, r, err)
			return
		}

		publicMediaID, err := utils.GenerateID()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if contentType == "" {
			contentType = "application/octet-stream"
		}

		_, err = app.S3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String("media"),
			Key:           aws.String(fileName),
			Body:          file,
			ContentType:   aws.String(contentType),
			ContentLength: &fileSize,
		})

		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		app.logger.Info("Successfully uploaded bytes: ", "filename", fileName, "size", fileSize)

		media, err := app.queries.CreateMedia(r.Context(), data.CreateMediaParams{
			UserID:        userID,
			PublicMediaID: publicMediaID,
			Filename:      fileName,
			MimeType:      contentType,
			Size:          fileSize,
		})

		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// w.WriteHeader(http.StatusCreated)
		// w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(`{"status":"success","filename":"` + fileName + `"}`))

		data := envelope{
			"status":   "success",
			"filename": fileName,
			"mediaid":  media.PublicMediaID,
		}

		err = app.writeJSON(w, http.StatusCreated, data, nil)

		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}

func (app *application) handleMediaGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//GET MEDIA BY MEDIA ID FROM URL PARAM
		// mediaID := r.URL.Query().Get("mediaid")

		// objectname, err := app.queries.GetMediaByPublicID(r.Context(), mediaID)
		// if err != nil {
		// 	app.serverErrorResponse(w, r, err)
		// 	return
		// }

		// app.presignClient.PresignGetObject(r.Context(), &s3.GetObjectInput{
		// 	Bucket: aws.String("media"),
		// 	Key:    aws.String(objectname),
		// })
	}
}

func (app *application) handleMediaList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//GET ALL MEDIA BY USER ID
	}
}

func (app *application) handleMediaDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
