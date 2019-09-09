package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rbobillo/OnDiraitDeLaMagie/families/dto"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

// Conn is the main connection to rabbit
var Conn *amqp.Connection

// Chan is the main rabbit channel
var Chan *amqp.Channel

// Pubq are all the queues
// where Hogwarts should publish in
var Pubq = make(map[string]amqp.Queue)

// Subq is the queue Hogwarts listens to
var Subq amqp.Queue

// DeclareBasicQueue is used to declare once a
// RabbitMQ queue, with default parameters
func DeclareBasicQueue(name string) amqp.Queue{
	q, err := Chan.QueueDeclare(name,
		false,
		false,
		false,
		false,
		nil,
	)
	HandleError(err, fmt.Sprintf("failed to declare the queue %s"), Error)
	return q
}

// Publish sends payload to 'pubq'
func Publish(qname string, payload string){
	err := Chan.Publish(
		"",
		Pubq[qname].Name,
		false,
		false,
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(payload),
		})
	HandleError(err, "failed to publish a message", Error)
}

// Subscribe listens to 'subq' (families)
// Each time a message is received
// it is parsed and handled
func Subscribe(w *http.ResponseWriter) {
	msgs, err := Chan.Consume(
		Subq.Name,			// queue
		"",		// consumer
		false,		// auto-ack (should the message be removed from queue after beind read)
		false,		// exclusive
		false,		// no-local
		false,		// no-wait
		nil,			// args
	)
	HandleError(err, "Failed to register a consumer", Warn)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a mail: %s", d.Body)

			// TODO: check message content, and publish on condition, to the right queue

			if d.Body != nil {

				var alert    dto.Alert
				var eligible dto.Eligible
				var safety   dto.Safety

				cannotParseAlert 	:= json.Unmarshal(d.Body, &alert)    // check if 'alert' is well created ?
				cannotParseEligible := json.Unmarshal(d.Body, &eligible) // check if 'eligible' is well created ?
				cannotParseSafety   := json.Unmarshal(d.Body, &safety)   // check if 'safety' is well created ?

				if cannotParseAlert == nil {
					//AlertHogwarts(alert)
					d.Ack(false)
				} else if cannotParseEligible == nil {
					AttendHogwarts(eligible)

					d.Ack(false)
					// try to parse another type of message, or fail
				} else if cannotParseSafety == nil {

				}
			}
		}
	}()

	log.Printf("Waiting for mails...")

	<-forever
}

// AttendHogwarts evaluates the emergency
// and helps Hogwarts
func AttendHogwarts(eligible dto.Eligible) {
	hogwartsURL := GetEnvOrElse("HOGWARTS_URL", "http://localhost:9091")

	attendEndpoint := "/actions/" + eligible.WizardID.String() + "/attend"

	eligibleWizard, err := json.Marshal(eligible)

	req, err := http.NewRequest("PATCH", hogwartsURL+attendEndpoint, bytes.NewBuffer(eligibleWizard))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}

//func AlertHogwarts(alert dto.Alert, Conn amqp.Connection) {
//	Conn.addB
//		addBlockedListener(new BlockedListener() {
//		public void handleBlocked(String reason) throws IOException {
//			// Connection is now blocked
//		}
//
//		public void handleUnblocked() throws IOException {
//			// Connection is now unblocked
//		}
//	});
//}