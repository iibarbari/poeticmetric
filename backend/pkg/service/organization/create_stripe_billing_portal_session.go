package organization

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/billingportal/session"
	"github.com/th0th/poeticmetric/backend/pkg/depot"
	"github.com/th0th/poeticmetric/backend/pkg/frontend"
	"github.com/th0th/poeticmetric/backend/pkg/model"
	"github.com/th0th/poeticmetric/backend/pkg/pointer"
)

func CreateStripeBillingPortalSession(dp *depot.Depot, id uint64) (*stripe.BillingPortalSession, error) {
	modelOrganization := &model.Organization{}

	err := dp.Postgres().
		Model(&model.Organization{}).
		Select("stripe_customer_id").
		Where("id = ?", id).
		First(modelOrganization).
		Error
	if err != nil {
		return nil, err
	}

	if modelOrganization.StripeCustomerId == nil {
		return nil, ErrNotStripeCustomer
	}

	stripeBillingPortalSession, err := session.New(&stripe.BillingPortalSessionParams{
		Customer:  modelOrganization.StripeCustomerId,
		ReturnURL: pointer.Get(frontend.GenerateUrl("/billing")),
	})
	if err != nil {
		return nil, err
	}

	return stripeBillingPortalSession, nil
}
