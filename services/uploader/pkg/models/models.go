package models

type FileMetaData struct {
	FileName string `json:file_name`
	FileSize int64  `json:file_size`
	URL      string `json:url`
}

