package main

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/skip2/go-qrcode"
)

type templateData struct {
	Title   string
	Message string
	QR      string
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "index.html")
	page := templateData{Title: "QR code generator"}
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		app.logger.Err(err).Msg("Error ParseFile homeHandler")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, page); err != nil {
		app.logger.Err(err).Msg("Error Execute homeHandler")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

type qrInputForm struct {
	Message string `form:"message"`
}

func (app *application) generateQRHandler(w http.ResponseWriter, r *http.Request) {
	var form qrInputForm
	err := r.ParseForm()
	if err != nil {
		app.logger.Err(err).Msg("Error Parse Form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {

		app.logger.Err(err).Msg("Error decode Form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	totalLength := utf8.RuneCountInString(form.Message)

	if totalLength == 0 || totalLength > 300 {
		app.logger.Err(errors.New("InvalidMessageInput")).Msg("Message must be between 0-300")
		http.Error(w, "Message must be between 0-300", http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprint(time.Now().Unix()) + ".png"
	err = qrcode.WriteFile(form.Message, qrcode.Medium, 256, "./static/"+fileName)

	if err != nil {
		app.logger.Err(err).Msg("Error generate QR code")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fp := path.Join("templates", "index.html")
	page := templateData{Title: "DCOM QR code generator", Message: form.Message, QR: "/static/" + fileName}
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		app.logger.Err(err).Msg("Error ParseFile homeHandler")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, page); err != nil {
		app.logger.Err(err).Msg("Error Execute homeHandler")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
