package cli

import (
	"boilerplate-api/lib/config"
)

// Migrate command for running migrations with gap filling
type Migrate struct {
	logger    config.Logger
	migration *config.Migrations
}

// NewMigrate creates a new migrate command
func NewMigrate(
	logger config.Logger,
	migration *config.Migrations,
) Migrate {
	return Migrate{
		logger:    logger,
		migration: migration,
	}
}

// Name returns the command name
func (c Migrate) Name() string {
	return "Migrate"
}

// Run executes the migrate up command
func (c Migrate) Run() {
	if err := c.migration.MigrateUp(); err != nil {
		c.logger.Error("Migration failed: ", err)
		return
	}
	c.logger.Info("✅ Migration completed successfully")
}

// MigrateUp is a helper for explicit up migration (used from CLI or main)
func (c Migrate) MigrateUp() error {
	return c.migration.MigrateUp()
}

// MigrateDown is a helper for explicit down migration
func (c Migrate) MigrateDown() error {
	return c.migration.MigrateDown()
}
