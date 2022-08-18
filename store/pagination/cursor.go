package pagination

import (
	"encoding/base64"
	"fmt"
)

func EncodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}

func decodeCursor(cursor string) (string, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", fmt.Errorf("invalid cursor")
	}

	return string(decodedCursor), nil
}
