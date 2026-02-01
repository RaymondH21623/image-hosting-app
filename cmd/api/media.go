package main

import (
	"net/http"
	"shareapp/internal/data"
	"shareapp/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
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
		mediaID := chi.URLParam(r, "id")
		app.logger.Info("Fetching media with ID: ", "mediaID", mediaID)

		objectname, err := app.queries.GetMediaNameByPublicID(r.Context(), mediaID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		presignResult, err := app.presignClient.PresignGetObject(r.Context(), &s3.GetObjectInput{
			Bucket: aws.String("media"),
			Key:    aws.String(objectname),
		})
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		app.logger.Info(presignResult.URL)

		data := envelope{
			"status": "success",
			"url":    presignResult.URL,
		}

		err = app.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *application) handleMediaListGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//GET ALL MEDIA BY USER ID
		userID := chi.URLParam(r, "id")
		app.logger.Info("Listing media for user ID: ", "userID", userID)

		mediaList, err := app.queries.ListMediaByUser(r.Context(), userID)

		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		data := envelope{
			"status": "success",
			"data":   mediaList,
		}

		err = app.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *application) handleMediaDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
