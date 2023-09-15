package helpers

import (
	"net/http"
	"social-api/logger"
)

func HandleParserError(parseError error, w http.ResponseWriter, log logger.Logger, msg ...string) {
	var message string
	if len(msg) == 1 && msg[0] != " " {
		message = msg[0]
	} else {
		message = ""
	}
	switch parseError.Error() {
	case "failed to readAll of byte stream":
		log.WriteToLogger(logger.ERROR, "couldnt read all of the byte stream:"+message)
		fallthrough
	case "error when unmarshaling the data into generic":
		log.WriteToLogger(logger.ERROR, "error when unmarshing data into generic:"+message)
		fallthrough
	default:
		log.WriteToLogger(logger.ERROR, "unknow error when parsing request", parseError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error when parsing the request"))
	}
}
