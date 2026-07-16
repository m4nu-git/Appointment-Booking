package domain

// Availability maps onto the `availability` table created in migration
// 000003. start_time/end_time are Postgres TIME columns; GORM/pgx will
// scan them into Go's time.Time with a zero date component. We keep them
// as strings ("HH:MM") at the domain boundary instead, and convert at the
// repository layer — this avoids leaking Postgres's TIME-as-time.Time quirk
// into business logic, and matches the HH:MM contract used everywhere else
// in the API (request/response DTOs, slotId format).
type Availability struct {
	ID         string `gorm:"column:id;primaryKey"`
	ServiceID  string `gorm:"column:service_id"`
	DayOfWeek  int    `gorm:"column:day_of_week"` // 0 = Sunday ... 6 = Saturday
	StartTime  string `gorm:"column:start_time"`  // "HH:MM", converted in repository
	EndTime    string `gorm:"column:end_time"`
}

func (Availability) TableName() string { return "availability" }