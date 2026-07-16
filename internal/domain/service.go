package domain

import "time"

type ServiceType string

const (
	ServiceTypeMedical   ServiceType = "MEDICAL"
	ServiceTypeHouseHelp ServiceType = "HOUSE_HELP"
	ServiceTypeBeauty    ServiceType = "BEAUTY"
	ServiceTypeFitness   ServiceType = "FITNESS"
	ServiceTypeEducation ServiceType = "EDUCATION"
	ServiceTypeOther     ServiceType = "OTHER"
)

// Service maps onto the `services` table created in migration 000002.
type Service struct {
	ID              string      `gorm:"column:id;primaryKey"`
	Name            string      `gorm:"column:name"`
	Type            ServiceType `gorm:"column:type"`
	ProviderID      string      `gorm:"column:provider_id"`
	DurationMinutes int         `gorm:"column:duration_minutes"`
	CreatedAt       time.Time   `gorm:"column:created_at;autoCreateTime"`
}

func (Service) TableName() string { return "services" }