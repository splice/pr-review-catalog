package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// parseJSON decodes r as JSON into dst, or returns an error.
func (pc *pageController) parseJSON(r io.Reader, dst interface{}) error {
	return json.NewDecoder(r).Decode(dst)
}

// renderJSON writes v as JSON into w, or returns an error.
func (pc *pageController) renderJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// serverError is invoked whenever a "500 Internal Server Error" response is needed.
func (pc *pageController) serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// badRequestError is invoked whenever a "400 Bad Request Error" response is needed.
func (pc *pageController) badRequestError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (pc *pageController) parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		// if parse form fails, it is because the form data is invalid
		return err
	}

	// add the path variables to the map - it should take precedence over
	// any query parameter.
	vars := mux.Vars(r)
	for k, v := range vars {
		r.Form.Set(k, v)
	}

	if err := pc.dec.Decode(dst, r.Form); err != nil {
		return err
	}
	return nil
}

func (pc *pageController) parseJSONForm(r io.Reader, dst interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		// this is an I/O error, do not wrap in a validation error
		return err
	}
	if len(b) > 0 {
		// if unmarshal fails, it is because the JSON is invalid
		return errors.Wrap(json.Unmarshal(b, dst), "json is invalid")
	}
	return nil
}
