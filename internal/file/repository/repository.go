package repository

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
	"gorm.io/gorm"
)

type RepositoryFile struct {
	db gorm.DB
}

func NewRepositoryFile(db gorm.DB) RepositoryFile {
	return RepositoryFile{db: db}
}

func (r RepositoryFile) GetFileByID(ctx context.Context, id string) (res entity.File, err error) {
	err = r.db.WithContext(ctx).First(&res, "id = ?", id).Error
	if err != nil {
		return entity.File{}, err
	}
	return res, nil
}
