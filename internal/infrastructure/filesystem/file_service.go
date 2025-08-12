package filesystem

import (
	"Blog-API/internal/domain"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fileService struct {
	uploadPath string
}

func NewFileService(uploadPath string) domain.FileService {
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		// Log a fatal error because the app can't run without this directory.
		panic(fmt.Sprintf("could not create upload directory: %v", err))
	}
	return &fileService{
		uploadPath: uploadPath,
	}
}

// SaveProfilePicture implement the interface method.
func (s *fileService) SaveProfilePicture(userID primitive.ObjectID, file multipart.File, handler *multipart.FileHeader) (*domain.Photo, error) {
	ext := filepath.Ext(handler.Filename)
	newFileName := fmt.Sprintf("%s-%d%s", userID.Hex(), time.Now().Unix(), ext)
	path := filepath.Join(s.uploadPath, newFileName)

	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	photo := &domain.Photo{
		Filename:   newFileName,
		FilePath:   path,
		PublicID:   "",
		UploadedAt: time.Now(),
	}

	return photo, nil
}
