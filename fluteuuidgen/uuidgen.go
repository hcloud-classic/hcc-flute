package cellouuidgen

import (
	"GraphQL_Cello/cellologger"
	"github.com/nu7hatch/gouuid"
)

// Uuidgen : Generate uuid
func Uuidgen() (string, error) {
	out, err := uuid.NewV4()
	if err != nil {
		cellologger.Logger.Println(err)
		return "", err
	}

	return out.String(), nil
}
