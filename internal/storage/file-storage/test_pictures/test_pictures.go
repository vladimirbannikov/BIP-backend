package test_pictures

import "path"

const pathPicture = "file-storage/tests_pictures"

type FileProvider interface {
	GetFile(path string) ([]byte, error)
}

type Repository struct {
	f FileProvider
}

func NewRepository(f FileProvider) *Repository {
	return &Repository{f: f}
}

func (r *Repository) GetPictureByte(pic string) ([]byte, error) {
	bytes, err := r.f.GetFile(path.Join(pathPicture, pic))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
