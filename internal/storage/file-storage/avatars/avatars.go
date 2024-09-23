package avatars

import (
	"math/rand"
	"path"
)

const avatarsPath = "file-storage/user-avatars"

type FileProvider interface {
	GetFile(path string) ([]byte, error)
	GetAllFileNames(path string) ([]string, error)
}

type Repository struct {
	f FileProvider
}

func NewRepository(f FileProvider) *Repository {
	return &Repository{f: f}
}

func (r *Repository) GetAvatarByte(avatar string) ([]byte, error) {
	bytes, err := r.f.GetFile(path.Join(avatarsPath, avatar))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (r *Repository) GetRandomAvatarFileName() (string, error) {
	names, err := r.f.GetAllFileNames(avatarsPath)
	if err != nil {
		return "", err
	}
	avatarName := names[rand.Intn(len(names))]
	return avatarName, nil
}
