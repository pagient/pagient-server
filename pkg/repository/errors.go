package repository

import (
	"os"
	"strings"
)

type entryNotExistErr struct {
	msg string
}

func (err *entryNotExistErr) Error() string {
	return err.msg
}

func (err *entryNotExistErr) EntryNotExist() bool {
	return true
}

type entryExistErr struct {
	msg string
}

func (err *entryExistErr) Error() string {
	return err.msg
}

func (err *entryExistErr) EntryExist() bool {
	return true
}

func isNotFoundErr(err error) bool {
	if strings.HasPrefix(err.Error(), "Unable to find file or directory") || os.IsNotExist(err) {
		return true
	}
	return false
}
