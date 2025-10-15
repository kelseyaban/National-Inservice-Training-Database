package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

// create an envelope type
type envelope map[string]any

func (a *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	jsResponse = append(jsResponse, '\n')
	// additional  headers to be set
	for key, value := range headers {
		w.Header()[key] = value
	}

	// set content type header
	w.Header().Set("Content-Type", "application/json")

	// explicitly set the response status code
	w.WriteHeader(status)
	_, err = w.Write(jsResponse)
	if err != nil {
		return err
	}

	return nil
}

// ensures that there is nothing can affect our backend
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	// max size of the request body(250KB reasonable)
	maxBytes := 256_000
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// decoder checking for unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Start the decoding
	err := dec.Decode(destination)

	if err != nil {
		// check for different errors
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// Decode can also send back an io error message
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("the body contains the incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("the body contains the incorrect  JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("the body must not be empty")
		// check for unknown field error
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknownfield")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		// does the body exceed our limit?
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("the body must not be larger than %d bytes", maxBytesError.Limit)
		// the programmer messed up
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})

	if !errors.Is(err, io.EOF) { // more data present
		return errors.New("the body must only contain a single JSON value")
	}

	return nil
}

// Getting the idfromt he URL
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// Get the URL parameters
	params := httprouter.ParamsFromContext(r.Context())

	// Convert the id from string to int
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) getSingleQueryParameter(queryParameters url.Values, key string, defaultValue string) string {
	// url.Values is a key:value hash map of the query parameters
	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	return result
}

// call when we have multiple comma-separated values
func (app *application) getMultipleQueryParameters(queryParameters url.Values, key string, defaultValue []string) []string {
	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	return strings.Split(result, ",")
}

// this method can cause a validation error when trying to convert the
// string to a valid integer value
func (app *application) getSingleIntegerParameter(queryParameters url.Values, key string, defaultValue int, v *validator.Validator) int {

	result := queryParameters.Get(key)
	if result == "" {
		return defaultValue
	}
	// try to convert to an integer
	intValue, err := strconv.Atoi(result)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return intValue
}

// Accept a function and run it in the background also recover from any panic
func (a *application) background(fn func()) {
	a.wg.Add(1) // Use a wait group to ensure all goroutines finish before we exit
	go func() {
		defer a.wg.Done() // signal goroutine is done
		defer func() {
			err := recover()
			if err != nil {
				a.logger.Error(fmt.Sprintf("%v", err))
			}
		}()
		fn() // Run the actual function
	}()
}
