package pathBuilder

import "navigation/internal/models"

type Repository interface {
	GetSectorLink() ([]models.SectorLink, error)
	GetSector(number string, building uint) (int, error)
}
