package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func AdminListUser() (map[string]any, error) {
	var adminListUserData []entities.AdminListUserResponse

	queryUsers := `SELECT p.user_id AS id, u.email, u.phone, p.avatar, p.fullname, ur.name AS role, u.created_at
		FROM users u 
		INNER JOIN profiles p ON p.user_id = u.uid
		INNER JOIN user_roles ur ON ur.id = u.role
	`

	rows, err := db.Debug().Raw(queryUsers).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var adminListUserQuery entities.AdminListUser

		if err := rows.Scan(
			&adminListUserQuery.Id,
			&adminListUserQuery.Email,
			&adminListUserQuery.Phone,
			&adminListUserQuery.Avatar,
			&adminListUserQuery.Fullname,
			&adminListUserQuery.Role,
			&adminListUserQuery.CreatedAt,
		); err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, errors.New(err.Error())
		}

		var adminListUserBranchData []entities.AdminListUserBranch
		queryBranches := `SELECT b.id, b.name
			FROM user_branches ub 
			LEFT JOIN branchs b ON b.id = ub.branch_id 
			WHERE ub.user_id = ?`
		if errUserBranch := db.Debug().Raw(queryBranches, adminListUserQuery.Id).Scan(&adminListUserBranchData).Error; errUserBranch != nil {
			helper.Logger("error", "In Server: "+errUserBranch.Error())
			return nil, errors.New(errUserBranch.Error())
		}

		var branch entities.AdminListUserBranch
		if len(adminListUserBranchData) > 0 {
			branch = adminListUserBranchData[0] // Get the first branch
		} else {
			branch = entities.AdminListUserBranch{
				Id:   0,
				Name: "-",
			}
		}

		adminListUserData = append(adminListUserData, entities.AdminListUserResponse{
			Id:        adminListUserQuery.Id,
			Avatar:    adminListUserQuery.Avatar.String,
			Fullname:  adminListUserQuery.Fullname,
			Email:     adminListUserQuery.Email,
			Phone:     adminListUserQuery.Phone,
			Role:      adminListUserQuery.Role,
			Branch:    branch,
			CreatedAt: adminListUserQuery.CreatedAt,
		})
	}

	return map[string]any{
		"data": adminListUserData,
	}, nil
}
