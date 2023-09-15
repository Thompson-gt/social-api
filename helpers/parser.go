package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"social-api/types"
)

// takes in io.ReaderCloser (request body) and unmarshals the request
// into the val (type bounded by Requesttypes in types package)
// returns a pointer to this newly filled reqeust Type (val should be a empty struct of any RequestType)
func ParseBody[T types.AuthUserRequest | types.RequestPost | types.RequestUser](body io.ReadCloser, val T) (*T, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, errors.New("failed to readAll of byte stream")
	}
	if err := json.Unmarshal(b, &val); err != nil {
		return nil, errors.New("error when unmarshaling the data into generic")
	}
	return &val, nil

}
