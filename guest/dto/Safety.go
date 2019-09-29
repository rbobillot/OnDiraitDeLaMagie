package dto

import "github.com/satori/go.uuid" // go get github.com/satori/go.uuid

// Safety is the message sent by Hogwarts (via Owls)
// to guest, to inform that Hogwarts is
// no longer under attack
type Safety struct {
	ID            uuid.UUID `json:"id"`
	WizardID      uuid.UUID `json:"wizard_id"`
	SafetyMessage string    `json:"safety_message"`
}