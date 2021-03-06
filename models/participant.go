package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Participant are the members in the local FDF group
type Participant struct {
	ID           uuid.UUID          `json:"id" db:"id"`
	UserID       uuid.UUID          `db:"user_id"`
	MemberNumber string             `json:"member_number" db:"member_number"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" db:"updated_at"`
	FirstName    string             `json:"first_name" db:"first_name"`
	LastName     string             `json:"last_name" db:"last_name"`
	Phone        string             `json:"phone" db:"phone"`
	DateOfBirth  nulls.Time         `json:"date_of_birth" db:"date_of_birth"`
	Classes      Classes            `many_to_many:"class_memberships" db:"-"`
	Memberships  ClassMemberships   `has_many:"class_memberships" db:"-"`
	Image        *ParticipantsImage `has_one:"participants_images" fk_id:"participant_id"`
}

// SelectLabel is used for creating dropdown boxes in plush
func (p Participant) SelectLabel() string {
	return p.FirstName
}

// SelectValue is used for creating dropdown boxes in plush
func (p Participant) SelectValue() interface{} {
	return p.ID
}

// String is not required by pop and may be deleted
func (p Participant) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Participants is not required by pop and may be deleted
type Participants []Participant

// String is not required by pop and may be deleted
func (p Participants) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Participant) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: p.LastName, Name: "LastName"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Participant) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Participant) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
