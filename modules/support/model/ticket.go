package model

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

// Model is the that all of the other models inherit.
type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        string         `gorm:"primaryKey"`
}

// Ticket is the model for a ticket.
type Ticket struct {
	Model
	ChannelID string
	SenderID  string
}

func (t *Ticket) BeforeCreate(*gorm.DB) error {
	t.ID = xid.New().String()

	return nil
}
