package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	entranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

// validate holds the validator settings and caches for validating request struct values.
var validate = validator.New()

// translator is a cache of locale and translation information.
var translator *ut.UniversalTranslator

func init() {

	// Init the english locale
	enLocale := en.New()

	// Create the translator with en as the fallback locale.
	translator = ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	_ = entranslations.RegisterDefaultTranslations(validate, lang)

	// Validate with JSON tag names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Params returns the web call parameters from the request.
func Params(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}

// Decode reads the HTTP request body and decode it into the provided value
// If the provided value is a struct then it is checked for validation tags.
// It also validates for field types.
func Decode(r *http.Request, val interface{}) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(val); err != nil {
		switch t := err.(type) {
		case *json.UnmarshalTypeError:
			fields := []FieldError{{fmt.Sprintf("cannot use field %s of type %s", t.Field, t.Value), t.Field}}
			return NewFieldsValidationError(fields)
		default:
			return NewGenericError(err, http.StatusBadRequest)
		}
	}

	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls the language of the error messages. You could look at the
		// Accept-Language header if you intend to support multiple languages.
		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}

		return NewFieldsValidationError(fields)
	}

	return nil
}
