package repository

type fileDriver interface {
	Write(string, string, interface{}) error
	Read(string, string, interface{}) error
	ReadAll(string string) ([]string, error)
	Delete(string, string) error
}

