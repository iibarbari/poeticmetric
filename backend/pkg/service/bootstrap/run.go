package bootstrap

import (
	"fmt"

	v "github.com/RussellLuo/validating/v3"
	"github.com/th0th/poeticmetric/backend/pkg/depot"
	"github.com/th0th/poeticmetric/backend/pkg/env"
	"github.com/th0th/poeticmetric/backend/pkg/model"
	"github.com/th0th/poeticmetric/backend/pkg/pointer"
	"github.com/th0th/poeticmetric/backend/pkg/service/userpassword"
	"github.com/th0th/poeticmetric/backend/pkg/service/userself"
)

type Payload struct {
	CreateDemoSite   bool    `json:"createDemoSite"`
	OrganizationName *string `json:"organizationName"`
	UserEmail        *string `json:"userEmail"`
	UserName         *string `json:"userName"`
	UserNewPassword  *string `json:"userNewPassword"`
	UserNewPassword2 *string `json:"userNewPassword2"`
}

func Run(dp *depot.Depot, payload *Payload) (*userself.UserSelf, error) {
	err := validateBootstrapPayload(payload)
	if err != nil {
		return nil, err
	}

	var modelPlans []*model.Plan

	if env.GetIsHosted() {
		modelPlans = []*model.Plan{
			{
				Id:                1,
				MaxEventsPerMonth: pointer.Get(uint64(100000)),
				MaxUsers:          pointer.Get(uint64(1)),
				Name:              "Basic",
				StripeProductId:   pointer.Get("prod_KXK6a9Zmy3qcLz"),
			},
			{
				Id:                2,
				MaxEventsPerMonth: pointer.Get(uint64(1000000)),
				MaxUsers:          pointer.Get(uint64(3)),
				Name:              "Pro",
				StripeProductId:   pointer.Get("prod_KXK7HFnQGBmP6D"),
			},
			{
				Id:                3,
				MaxEventsPerMonth: pointer.Get(uint64(5000000)),
				MaxUsers:          pointer.Get(uint64(50)),
				Name:              "Business",
				StripeProductId:   pointer.Get("prod_KXK83fu8EQrKfM"),
			},
		}
	} else {
		modelPlans = []*model.Plan{
			{
				Id:   1,
				Name: "Default",
			},
		}
	}

	modelOrganization := &model.Organization{
		Id:        1,
		IsOnTrial: false,
		Name:      *payload.OrganizationName,
		PlanId:    &modelPlans[len(modelPlans)-1].Id,
	}

	userPasswordHash, err := userpassword.GetHash(*payload.UserNewPassword)
	if err != nil {
		return nil, err
	}

	modelUser := &model.User{
		Email:               *payload.UserEmail,
		Id:                  1,
		IsActive:            true,
		IsEmailVerified:     true,
		IsOrganizationOwner: true,
		Name:                *payload.UserName,
		OrganizationId:      modelOrganization.Id,
		Password:            userPasswordHash,
	}

	err = dp.WithPostgresTransaction(func(dp2 *depot.Depot) error {
		err2 := dp2.Postgres().
			Create(modelPlans).
			Error
		if err2 != nil {
			return err2
		}

		err2 = dp2.Postgres().
			Create(modelOrganization).
			Error
		if err2 != nil {
			return err2
		}

		err2 = dp2.Postgres().
			Create(modelUser).
			Error
		if err2 != nil {
			return err2
		}

		if payload.CreateDemoSite {
			err2 = createSite(dp2)
			if err2 != nil {
				return err2
			}
		}

		var sequenceFixingQueries []string

		err2 = dp.Postgres().Raw(`
			SELECT 'SELECT SETVAL('||
				   "quote_literal"("quote_ident"("PGT"."schemaname") || '.' || "quote_ident"("S"."relname")) ||
				   ', COALESCE(MAX(' || "quote_ident"("C"."attname")|| '), 1) ) FROM ' ||
				   "quote_ident"("PGT"."schemaname")|| '.'||"quote_ident"("T"."relname")|| ';'
			FROM "pg_class" AS "S",
				 "pg_depend" AS "D",
				 "pg_class" AS "T",
				 "pg_attribute" AS "C",
				 "pg_tables" AS "PGT"
			WHERE "S"."relkind" = 'S'
				AND "S"."oid" = "D"."objid"
				AND "D"."refobjid" = "T"."oid"
				AND "D"."refobjid" = "C"."attrelid"
				AND "D"."refobjsubid" = "C"."attnum"
				AND "T"."relname" = "PGT"."tablename"
			ORDER BY "S"."relname";
		`).
			Scan(&sequenceFixingQueries).
			Error
		if err != nil {
			return err
		}

		for _, sequenceFixingQuery := range sequenceFixingQueries {
			err2 = dp.Postgres().
				Exec(sequenceFixingQuery).
				Error
			if err2 != nil {
				return err2
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userself.Read(dp, modelUser.Id)
}

func validateBootstrapPayload(payload *Payload) error {
	errs := v.Validate(v.Schema{
		v.F("organizationName", payload.OrganizationName): v.All(
			v.Nonzero[*string]().Msg("This field is required."),

			v.Is(func(t *string) bool {
				return len(*t) >= model.OrganizationNameMinLength && len(*t) <= model.OrganizationNameMaxLength
			}).Msg(fmt.Sprintf(
				"This field should be between %d and %d characters in length.",
				model.OrganizationNameMinLength,
				model.OrganizationNameMaxLength,
			)),
		),

		v.F("userEmail", payload.OrganizationName): v.All(
			v.Nonzero[*string]().Msg("This field is required."),
		),

		v.F("userName", payload.UserName): v.All(
			v.Nonzero[*string]().Msg("This field is required."),

			v.Is(func(t *string) bool {
				return len(*t) >= model.UserNameMinLength && len(*t) <= model.UserNameMaxLength
			}).Msg(fmt.Sprintf(
				"This field should be between %d and %d characters in length.",
				model.UserNameMinLength,
				model.UserNameMaxLength,
			)),
		),

		v.F("userNewPassword", payload.UserNewPassword): v.All(
			v.Nonzero[*string]().Msg("This field is required."),

			v.Is(func(t *string) bool {
				return len(*t) >= model.UserPasswordMinLength && len(*t) <= model.UserPasswordMaxLength
			}).Msg(fmt.Sprintf(
				"This field should be between %d and %d characters in length.",
				model.UserPasswordMinLength,
				model.UserPasswordMaxLength,
			)),
		),

		v.F("userNewPassword2", payload.UserNewPassword2): v.All(
			v.Nonzero[*string]().Msg("This field is required."),

			v.Is(
				func(t *string) bool {
					if payload.UserNewPassword == nil {
						return true
					}

					return *t == *payload.UserNewPassword
				}).Msg("Passwords don't match."),
		),
	})

	if len(errs) > 0 {
		return errs
	}

	return nil
}
