package domain

import "time"

type AppointmentStatus string

const (
	StatusBooked    AppointmentStatus = "BOOKED"
	StatusCancelled AppointmentStatus = "CANCELLED"
)

// Appointment maps onto the `appointments` table created in migration
// 000004. SlotID carries the unique constraint that guarantees no
// double-booking under concurrent requests (see migration comments).
type Appointment struct {
	ID        string            `gorm:"column:id;primaryKey"`
	UserID    string            `gorm:"column:user_id"`
	ServiceID string            `gorm:"column:service_id"`
	Date      string            `gorm:"column:date"` // "YYYY-MM-DD"
	StartTime string            `gorm:"column:start_time"`
	EndTime   string            `gorm:"column:end_time"`
	SlotID    string            `gorm:"column:slot_id"`
	Status    AppointmentStatus `gorm:"column:status"`
	CreatedAt time.Time         `gorm:"column:created_at;autoCreateTime"`
}

func (Appointment) TableName() string { return "appointments" }