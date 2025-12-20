package main

import (
	"io"
	"net/http"
	"shareapp/internal/data"
)

func (app *application) handleMediaPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		file, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		defer file.Close()

		fileName := fileHeader.Filename
		fileSize := fileHeader.Size
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		contentType := http.DetectContentType(fileBytes)

		// info, err := app.minio.PutObject(r.Context(), "media", fileName, bytes.NewReader(fileBytes), fileSize, minio.PutObjectOptions{ContentType: contentType})

		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		app.logger.Info("Successfully uploaded bytes: " + string(info.Size))

		media, err := app.queries.CreateMedia(r.Context(), data.CreateMediaParams{
			//UserID: ,
			Filename: fileName,
			MimeType: contentType,
			Size:     fileSize,
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
			"mediaid":  media.ID,
		}

		err = app.writeJSON(w, http.StatusCreated, data, nil)

		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
