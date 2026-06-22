package app

import (
	"log"
	"net/http"
)

func (a *App) decodeAndValidateForm(dst any, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		return err
	}

	if err := a.decoder.Decode(dst, r.PostForm); err != nil {
		log.Printf("Error decoding form: %v", err)
		return err
	}

	if err := a.validator.StructCtx(r.Context(), dst); err != nil {
		log.Printf("Error decoding form: %v", err)
		return err
	}

	return nil
}
