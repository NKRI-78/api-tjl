package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func BannerList() (map[string]interface{}, error) {
	banners := []entities.Banner{}

	query := `SELECT id, path, link FROM banners`

	err := db.Debug().Raw(query).Scan(&banners).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": banners,
	}, nil
}

func BannerStore(f *models.Banner) (map[string]interface{}, error) {
	banner := entities.Banner{}

	banner.Path = f.Path
	banner.Link = f.Link

	errInsertBanner := db.Debug().Exec(`
	INSERT INTO banners (path, link) 
	VALUES (?, ?)`, banner.Path, banner.Link).Error

	if errInsertBanner != nil {
		helper.Logger("error", "In Server: "+errInsertBanner.Error())
		return nil, errors.New(errInsertBanner.Error())
	}

	return map[string]any{}, nil
}
