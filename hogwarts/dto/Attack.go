package dto

import "github.com/satori/go.uuid"

// Attack is the action of
// the Villains attacking Hogwarts
// quick/strong determine how efficient
// the attack will be
type Attack struct {
	ID     uuid.UUID `json:"id"`
	Quick  bool      `json:"quick"`
	Strong bool      `json:"strong"`
}