package music

import (
	"path"
)

const musicPath = "file-storage/music"

type FileProvider interface {
	GetFile(path string) ([]byte, error)
}

type Repository struct {
	f FileProvider
}

func NewRepository(f FileProvider) *Repository {
	return &Repository{f: f}
}

func (r *Repository) GetMusicByte(music string) ([]byte, error) {
	bytes, err := r.f.GetFile(path.Join(musicPath, music))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
