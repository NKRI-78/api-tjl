package models

import (
	"superapps/entities"
)

type RegisterUserBranch entities.RegisterUserBranch
type UpdateEmail entities.UpdateEmail
type UpdateUser entities.UpdateUser

type User entities.User
type UserRoles entities.UserRoles
type UserAdmin entities.UserAdmin
type UserLogin entities.UserLogin
type UserOtp entities.UserOtp
type UserDelete entities.UserDelete

type ForgotPassword entities.ForgotPassword

type AdminListUser entities.AdminListUser
type AdminListUserResponse entities.AdminListUserResponse
type AdminListUserBranch entities.AdminListUserBranch

type GetBranch entities.Branch
type UpdateUserBranch entities.UpdateUserBranch

type Biodata entities.Biodata
type CheckAccount entities.CheckAccount

type Profile entities.Profile
type ProfileResponse entities.ProfileResponse
