package helpers

import (
	"errors"
	"net/http"
	"social-api/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

// will handle the error returned from the Modler interface
// handles the empty document as well as internal server errors
// optional can pass a message to be sent to the client
func HandleDbError(dbError error, w http.ResponseWriter, log logger.Logger, msg ...string) {
	if errors.Is(dbError, mongo.ErrNoDocuments) {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("item not found in the database"))
		return
	} else if len(msg) == 1 && msg[0] != "" {
		log.WriteToLogger(logger.ERROR, msg[0], dbError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg[0]))
		return
	} else {
		log.WriteToLogger(logger.ERROR, "unknown server error", dbError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unknow server error"))
		return
	}
}
