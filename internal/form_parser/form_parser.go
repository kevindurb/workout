package formparser

import (
	"log"
	"net/http"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
)

type FormParser struct {
	decoder   *form.Decoder
	validator *validator.Validate
}

func New() *FormParser {
	return &FormParser{
		decoder:   form.NewDecoder(),
		validator: validator.New(),
	}
}

func (fp *FormParser) Parse(dst any, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		return err
	}

	if err := fp.decoder.Decode(dst, r.PostForm); err != nil {
		log.Printf("Error decoding form: %v", err)
		return err
	}

	if err := fp.validator.StructCtx(r.Context(), dst); err != nil {
		log.Printf("Error validating form: %v", err)
		return err
	}

	return nil
}
