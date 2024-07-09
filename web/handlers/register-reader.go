package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"librarymanagement/db"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"net/http"
)

func RegisterReader(w http.ResponseWriter, r *http.Request) {

	var newReader db.Reader
	if err := json.NewDecoder(r.Body).Decode(&newReader); err != nil {
		slog.Error("Failed to decode new user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newReader,
		}))
		utils.SendError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	if err := utils.ValidateStruct(newReader); err != nil {
		slog.Error("Failed to validate new user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newReader,
		}))
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	newReader.Password = hashPassword(newReader.Password)

	var insertedReader *db.Reader
	var err error

	if insertedReader, err = db.GetReaderRepo().RegisterUser(&newReader); err != nil {
		utils.SendError(w, http.StatusExpectationFailed, err.Error())
		return
	}

	utils.SendData(w, insertedReader)

}

func hashPassword(pass string) string {

	h := sha256.New()
	h.Write([]byte(pass))
	hashValue := h.Sum(nil)
	return hex.EncodeToString(hashValue)
}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
