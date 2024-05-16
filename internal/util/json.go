package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Validator interface {
	Valid(ctx context.Context) (problems map[string]string)
}

var ErrValidation error
var ErrDecode error

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.MarshalIndent(payload, "", "   ")
	if err != nil {
		log.Println("Failed to marshal JSON response: ", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}

func DecodeJSON[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		ErrDecode = fmt.Errorf("error decoding JSON: %w", err)
		return v, nil, ErrDecode
	}

	if problems := v.Valid(r.Context()); len(problems) != 0 {
		ErrValidation = fmt.Errorf("invalid %T: %d problem(s)", v, len(problems))
		return v, problems, ErrValidation
	}

	return v, nil, nil
}
