package database

import "context"

// HealthCheck is an implementation of the HealthCheck interface,
// used to check the status of an app by ensuring it can connect to a database.
type HealthCheck struct {
	db Database
}

// NewHealthCheck returns a new instance of HealthCheck.
func NewHealthCheck(db Database) *HealthCheck {
	return &HealthCheck{db: db}
}

// Check ensures there is a successful connection to the database.
func (hc *HealthCheck) Check(ctx context.Context) error {
	return hc.db.Ping(ctx)
}
