package models

type PhotoUpload struct {
	Name string
	Data []byte
}

type CreateRequest struct {
	Fields map[string]string
	Photos []PhotoUpload
}
