package impl

import (
	"store-management-be/application/_dummp/validationapi"
	"store-management-be/application/auth/validation"
)

func init() {
	validationapi.UserInfo = validation.UserInfo
}
