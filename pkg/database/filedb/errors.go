package filedb

import (
	"strings"
	"os"
)

func isNotFoundErr(err error) bool {
	if strings.HasPrefix(err.Error(), "Unable to find file or directory") || os.IsNotExist(err) {
		return true
	}
	return false
}
