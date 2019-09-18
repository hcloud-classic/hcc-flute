package uuidgen

import (
	"hcloud-flute/logger"
	"github.com/nu7hatch/gouuid"
)

// Uuidgen : Generate uuid
func Uuidgen() (string, error) {
	out, err := uuid.NewV4()
	if err != nil {
		logger.Logger.Println(err)
		return "", err
	}

	return out.String(), nil
}
