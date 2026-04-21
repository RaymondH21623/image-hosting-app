package main

import (
	"net/http"
	"shareapp/internal/data"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (app *application) handleCreateMedia(w http.ResponseWriter, r *http.Request) {
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

	user := app.contextGetUser(r)
	userID := user.ID

	publicMediaID, err := app.generateNanoID()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = app.S3Storage.Client.PutObject(ctx, &s3.PutObjectInput{
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

	media := &data.Media{
		UserID:   userID,
		PublicID: publicMediaID,
		Filename: fileName,
		MimeType: contentType,
		Size:     fileSize,
	}

	err = app.models.Media.CreateMedia(media)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"status":   "success",
		"filename": fileName,
		"mediaid":  media.PublicID,
	}

	err = app.writeJSON(w, http.StatusCreated, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) handleShowMedia(w http.ResponseWriter, r *http.Request) {

	mediaID, err := app.readIDParam(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.Info("Fetching media with ID: ", "mediaID", mediaID)

	media, err := app.models.Media.GetMediaByPublicID(mediaID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	presignResult, err := app.S3Storage.PresignClient.PresignGetObject(r.Context(), &s3.GetObjectInput{
		Bucket: aws.String("media"),
		Key:    aws.String(media.Filename),
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

func (app *application) handleListMedia(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserID string
		Title  string
		data.Filters
	}

	userID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.Info("Listing media for user ID: ", "userID", userID)

	mediaList, err := app.models.Media.ListMediaByUser(userID)

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

func (app *application) handleDeleteMedia(w http.ResponseWriter, r *http.Request) {

}

func (app *application) handleUpdateMedia(w http.ResponseWriter, r *http.Request) {

}
