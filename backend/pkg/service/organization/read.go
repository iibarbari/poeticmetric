package organization

import (
	"github.com/poeticmetric/poeticmetric/backend/pkg/depot"
	"github.com/poeticmetric/poeticmetric/backend/pkg/model"
)

func Read(dp *depot.Depot, organizationId uint64) (*Organization, error) {
	organization := &Organization{}

	err := dp.Postgres().
		Model(&model.Organization{}).
		Joins("inner join plans on plans.id = organizations.plan_id").
		Select(
			"organizations.created_at",
			"organizations.id",
			"organizations.is_on_trial",
			"organizations.name",
			"organizations.stripe_customer_id",
			"organizations.subscription_period",
			"organizations.trial_ends_at",
			"organizations.updated_at",
			"plans.name as plan__name",
		).
		Where("organizations.id = ?", organizationId).
		First(organization).
		Error
	if err != nil {
		return nil, err
	}

	return organization, nil
}
