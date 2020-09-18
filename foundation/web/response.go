package web

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Respond sends a json data response converted from Go interface to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	// If the context object is missing the request values
	// do a graceful shutdown.
	rvs, ok := ctx.Value(KeyRequestValues).(*RequestValues)
	if !ok {
		return NewShutdownError("request values missing from context")
	}
	rvs.StatusCode = statusCode

	w.WriteHeader(statusCode)

	// If there is no data to return back to the client then return.
	if statusCode == http.StatusNoContent {
		return nil
	}

	// Convert the data from the interface to json.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Once we know that the marshall succeeded, next step is to set the content type to json and status code.
	w.Header().Set("Content-Type", "application/json")

	// Write the response to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func RespondWithRedirect(ctx context.Context, w http.ResponseWriter, r *http.Request, url string) error {
	statusCode := http.StatusMovedPermanently

	// If the context object is missing the request values
	// do a graceful shutdown.
	rvs, ok := ctx.Value(KeyRequestValues).(*RequestValues)
	if !ok {
		return NewShutdownError("request values missing from context")
	}
	rvs.StatusCode = statusCode

	http.Redirect(w, r, url, statusCode)

	return nil
}

func RespondWithError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error {

	// Check what kind of error we have and respond with the correct format
	switch e := errors.Cause(err).(type) {
	case *GenericError:
		er := GenericErrorResponse{Error: e.Err.Error()}
		return Respond(ctx, w, er, e.StatusCode)

	case *FieldsValidationError:
		er := FieldValidationErrorResponse{
			Error:  e.Err.Error(),
			Fields: e.Fields,
		}
		return Respond(ctx, w, er, e.StatusCode)

	case *RedirectError:
		http.Redirect(w, r, e.Url, e.StatusCode)
		return nil

	default:
		// If not, the handler send any arbitrary error so use 500 status
		return Respond(ctx, w, struct{}{}, http.StatusInternalServerError)
	}
}
