package file

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
)

type IRepositoryFile interface {
	GetFileByID(ctx context.Context, id string) (res entity.File, err error)
}
