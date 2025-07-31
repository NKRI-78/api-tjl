package services

import (
	"errors"

	entities "superapps/entities"
	helper "superapps/helpers"
)

func Branch() (map[string]any, error) {
	branches := []entities.Branch{}

	query := `SELECT id, name
	FROM branchs`

	err := db.Debug().Raw(query).Scan(&branches).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": branches,
	}, nil
}

func StoreBranch(sb *entities.StoreBranch) (map[string]any, error) {

	query := `INSERT INTO branchs (name) VALUES (?)`
	err := db.Debug().Exec(query, sb.Name).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func DeleteBranch(dbr *entities.DeleteBranch) (map[string]any, error) {
	query := `DELETE FROM branchs WHERE id = ?`
	err := db.Debug().Exec(query, dbr.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func UpdateBranch(ubr *entities.UpdateBranch) (map[string]any, error) {
	query := `UPDATE branchs SET name = ? WHERE id = ?`

	err := db.Debug().Exec(query, ubr.Name, ubr.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func AssignBranch(ab *entities.AssignBranch) (map[string]any, error) {

	query := `INSERT INTO user_branches(branch_id, user_id) VALUES(?, ?)`

	err := db.Debug().Exec(query, ab.BranchId, ab.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
