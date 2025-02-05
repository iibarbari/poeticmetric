package country

import (
	"github.com/th0th/poeticmetric/backend/pkg/country"
	"github.com/th0th/poeticmetric/backend/pkg/depot"
	"github.com/th0th/poeticmetric/backend/pkg/service/sitereport/filter"
	"gorm.io/gorm"
)

type Datum struct {
	Country           string `json:"country"`
	CountryIsoCode    string `json:"countryIsoCode"`
	VisitorCount      uint64 `json:"visitorCount"`
	VisitorPercentage uint16 `json:"visitorPercentage"`
}

type Report []*Datum

func Get(dp *depot.Depot, filters *filter.Filters) (Report, error) {
	report := Report{}

	baseQuery := filter.Apply(dp, filters).
		Where("country_iso_code is not null")

	totalVisitorCountSubQuery := baseQuery.
		Session(&gorm.Session{}).
		Select("count(distinct visitor_id) as count")

	err := baseQuery.
		Session(&gorm.Session{}).
		Joins("cross join (?) total_visitors", totalVisitorCountSubQuery).
		Select(
			"country_iso_code",
			"count(distinct visitor_id) as visitor_count",
			"toUInt16(round(100 * visitor_count / total_visitors.count)) as visitor_percentage",
		).
		Group("country_iso_code, total_visitors.count").
		Order("visitor_count desc, country_iso_code").
		Find(&report).
		Error
	if err != nil {
		return nil, err
	}

	for i := range report {
		report[i].Country = country.GetNameFromIsoCode(report[i].CountryIsoCode)
	}

	return report, nil
}
