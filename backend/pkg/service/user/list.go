package user

import (
	"github.com/poeticmetric/poeticmetric/backend/pkg/depot"
	"github.com/poeticmetric/poeticmetric/backend/pkg/model"
)

func List(dp *depot.Depot, organizationId uint64) ([]*User, error) {
	users := []*User{}

	err := dp.Postgres().
		Model(&model.User{}).
		Where("organization_id = ?", organizationId).
		Find(&users).
		Error

	return users, err
}
