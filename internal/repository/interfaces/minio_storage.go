package interfaces

type FileStorage interface {
	Upload(filename string, data []byte, contentType string) (string, error)
	Delete(filename string) error
	DeleteByURL(url string) error
	GetURL(filename string) string
}
