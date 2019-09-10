package fluteuuidgen

import (
	"GraphQL_Flute/flutelogger"
	"github.com/nu7hatch/gouuid"
)

// Uuidgen : Generate uuid
func Uuidgen() (string, error) {
	out, err := uuid.NewV4()
	if err != nil {
		flutelogger.Logger.Println(err)
		return "", err
	}

	return out.String(), nil
}
