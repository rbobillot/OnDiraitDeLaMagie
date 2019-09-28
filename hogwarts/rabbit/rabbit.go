package rabbit

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rbobillo/OnDiraitDeLaMagie/hogwarts/dto"
	"github.com/rbobillo/OnDiraitDeLaMagie/hogwarts/hogwartsinventory"
	"github.com/rbobillo/OnDiraitDeLaMagie/hogwarts/internal"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"log"
)

// Conn is the main connection to rabbit
var Conn *amqp.Connection

// Chan is the main rabbit channel
var Chan *amqp.Channel

// Pubq are all the queues
// where Ministry should publish in
var Pubq = make(map[string]amqp.Queue)

// Subq is the queue Ministry listens to
var Subq amqp.Queue

// Publish sends messages to 'pubq'
func Publish(qname string, payload string) {

	err := Chan.Publish(
		"",    		// exchange
		Pubq[qname].Name,			// routing key
		false,			// mandatory
		false,			// immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(payload),
		})

	internal.HandleError(err, "Failed to publish a message", internal.Warn)
}

// Subscribe listens to 'subq' (ministry)
// Each time a message is received
// it is parsed and handled
func Subscribe(db *sql.DB) {
	msgs, err := Chan.Consume(
		Subq.Name,			// queue
		"",		// consumer
		false,		// auto-ack (should the message be removed from queue after beind read)
		false,		// exclusive
		false,		// no-local
		false,		// no-wait
		nil,			// args
	)
	internal.HandleError(err, "Failed to register a consumer", internal.Warn)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a mail: %s", d.Body)

			// TODO: check message content, and publish on condition, to the right queue

			if d.Body != nil {

				var slot dto.Slot
				var arrested dto.Arrested
				var born dto.Birth


				dec := json.NewDecoder(bytes.NewReader(d.Body))
				dec.DisallowUnknownFields()
				cannotParseSlot := dec.Decode(&slot)

				dec = json.NewDecoder(bytes.NewReader(d.Body))
				dec.DisallowUnknownFields()
				cannotParseArrested := dec.Decode(&arrested)

				dec = json.NewDecoder(bytes.NewReader(d.Body))
				dec.DisallowUnknownFields()
				cannotParseBorn := dec.Decode(&born)


				if cannotParseSlot == nil {

					err, availableSlot := checkSlot(slot, db)
					if err != nil {
						internal.Warn(fmt.Sprintf("%s", err))

						err := d.Nack(true, true)
						if  err != nil {
							internal.Warn(fmt.Sprintf("cannot n.ack current message %s", slot.ID))
							return
						}
					}

					err = d.Ack(false)
					if err != nil {
						internal.Warn(fmt.Sprintf("cannot ack the current message : %s", slot.ID))
						return
					}
					err := slotHogwarts(availableSlot)


				} else  if cannotParseArrested == nil {

					err = d.Ack(false)
					if err != nil {
						internal.Warn(fmt.Sprintf("cannot ack the current message : %s", arrested.ID))
						return
					}

					internal.Debug("inform Guest and Families that Hogwarts is no longer under attack")

					safety, err := json.Marshal(dto.Safety{
						ID: 			uuid.Must(uuid.NewV4()),
						WizardID: 		arrested.WizardID,
						SafetyMessage: 		"Hogwarts is ready to receive new visits",
					})
					if err != nil {
						internal.Warn("cannot serialize Attack to JSON")
						return
					}

					Publish("families", string(safety))
					internal.Debug("Mail (safety) sent to families") //TODO: better message

					Publish("guest", string(safety))
					internal.Debug("Mail (safety) sent to guest") //TODO: better message


				}
			}
		}
	}()

	log.Printf("Waiting for mails...")

	<-forever
}

// DeclareBasicQueue is used to declare once
// a RabbitMQ queue, with default parameters
func DeclareBasicQueue(name string) amqp.Queue {
	q, err := Chan.QueueDeclare(name,
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	internal.HandleError(err, "Failed to declare a queue", internal.Warn)

	return q
}


func checkSlot(slot dto.Slot, db *sql.DB) (err error, available int ){

	query := "SELECT * FROM actions WHERE status = 'ongoing' and action = 'visit'"

	ongoing, err := hogwartsinventory.GetActions(db, query)
	if err !=  nil {
		internal.Warn("cannot get actions in hogwarts inventory")
		return err, 0
	}
	if len(ongoing) > 10 {
		return fmt.Errorf("hogwarts have 10 visit ongoing"), 0
	}
	return err, 9
}

func slotHogwarts(availableSlot int) (err error){
	available, err := json.Marshal(dto.Available{
		ID: 			uuid.Must(uuid.NewV4()),
		AvailableSlot:  availableSlot,
		AvailableMessage: 		"Hogwarts is ready to receive new visits",
	})
	Publish("guest", string(available))
	if err != nil {

	}
	return nil
}