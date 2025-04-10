package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func BannerList() (map[string]any, error) {
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

func BannerStore(f *models.Banner) (map[string]any, error) {
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

func BannerUpdate(f *models.Banner) (map[string]any, error) {
	banner := entities.Banner{}

	banner.Id = f.Id
	banner.Path = f.Path
	banner.Link = f.Link

	errUpdateBanner := db.Debug().Exec(`
		UPDATE banners SET path = ?, link = ?
		WHERE id = ?
	`, banner.Path, banner.Link, banner.Id).Error

	if errUpdateBanner != nil {
		helper.Logger("error", "In Server: "+errUpdateBanner.Error())
		return nil, errors.New(errUpdateBanner.Error())
	}

	return map[string]any{}, nil
}

func BannerDelete(f *models.Banner) (map[string]any, error) {
	banner := entities.Banner{}

	banner.Id = f.Id

	errDeleteBanner := db.Debug().Exec(`
		DELETE FROM banners WHERE id = ?
	`, banner.Id).Error

	if errDeleteBanner != nil {
		helper.Logger("error", "In Server: "+errDeleteBanner.Error())
		return nil, errors.New(errDeleteBanner.Error())
	}

	return map[string]any{}, nil
}
