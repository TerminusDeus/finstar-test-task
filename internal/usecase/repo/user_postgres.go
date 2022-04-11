package repo

import (
	"finstar-test-task/pkg/postgres"
	"finstar-test-task/proto"

	"gorm.io/gorm"
)

// UserRepo to map
type UserRepo struct {
	*postgres.Postgres
}

func (ur *UserRepo) IncreaseBalance(userId uint64, receipt float64) error {
	tx := ur.GetDatabase().Begin()
	if err := tx.Error; err != nil {
		return err
	}

	update := tx.Model(&proto.UserORM{}).Where("id = ?", userId).UpdateColumn("balance", gorm.Expr("balance  + ?", receipt))
	if update.Error != nil {
		tx.Rollback()

		return update.Error
	}

	if update.RowsAffected == 0 {
		tx.Rollback()

		return gorm.ErrRecordNotFound
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()

		return err
	}

	return nil
}

func (ur *UserRepo) TransferBalance(userIdFrom, userIdTo uint64, writeOff float64) error {
	tx := ur.GetDatabase().Begin()
	if err := tx.Error; err != nil {
		return err
	}

	update := tx.Model(&proto.UserORM{}).Where("id = ?", userIdFrom).UpdateColumn("balance", gorm.Expr("balance - ?", writeOff))
	if update.Error != nil {
		tx.Rollback()

		return update.Error
	}

	if update.RowsAffected == 0 {
		tx.Rollback()

		return gorm.ErrRecordNotFound
	}

	update = tx.Model(&proto.UserORM{}).Where("id = ?", userIdTo).UpdateColumn("balance", gorm.Expr("balance + ?", writeOff))
	if update.Error != nil {
		tx.Rollback()

		return update.Error
	}

	if update.RowsAffected == 0 {
		tx.Rollback()

		return gorm.ErrRecordNotFound
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()

		return err
	}

	return nil
}

// New -.
func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}
