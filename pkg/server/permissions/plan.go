package permissions

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// CheckPlanAllowance checks if the given user is within the allowance for the
// current plan.
func CheckPlanAllowance(db *gorm.DB, user database.User) (bool, error) {
	if user.Cloud {
		return true, nil
	}

	var bookCount int
	if err := db.Model(database.Book{}).Where("user_id = ? AND NOT deleted", user.ID).Count(&bookCount).Error; err != nil {
		return false, errors.Wrap(err, "checking plan threshold")
	}

	if bookCount >= 5 {
		return false, nil
	}

	return true, nil
}
