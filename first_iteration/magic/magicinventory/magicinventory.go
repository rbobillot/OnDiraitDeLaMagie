// Package magicinventory is used to setup and manipulate
// the magic database (magicinventory)
package magicinventory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rbobillo/OnDiraitDeLaMagie/first_iteration/magic/dao"
	"github.com/rbobillo/OnDiraitDeLaMagie/first_iteration/magic/internal"
	"log"
	"net/http"
)

//func getId(r *http.Request) string{
//	return mux.Vars(r)["id"]
//}
//type args interface {
//	getId() string
//}

// CreateWizard inserts a new Wizard into magicinventory
func CreateWizard(w dao.Wizard, db *sql.DB) (err error) {
	populateQuery :=
		`insert into wizards (id, first_name, last_name, age, category, arrested, dead)
                values ($1, $2, $3, $4, $5, $6, $7);`

	_, err = db.Exec(populateQuery, w.ID, w.FirstName, w.LastName, w.Age, w.Category, w.Arrested, w.Dead)

	if err != nil {
		log.Println("Cannot create wizard:", w, err)
		return err
	}

	log.Println("Created wizard:", w)

	return nil
}
// CheckWizardStatusById check if the status of an wizard is already set {jail, die, whatever)
//func CheckWizardStatusById(db *sql.DB, status string, id string)(err error){
//
//}
// GetWizardStatusById return the status of an wizard.
//func GetWizardStatusById(db *sql.DB, status string, id string)(err error){
//	_, err = db.Exec("SELECT FROM ")
//}

// UpdateWizard should update a Wizard in magicinventory
// TODO: must take a query
func UpdateWizards(db *sql.DB, status string, value float64) (err error) {

	_, err = db.Exec(fmt.Sprintf("UPDATE wizards SET %s = %s + $1", status, status), value)

	if err != nil {
		log.Println(fmt.Sprintf("Cannot update wizards %s", status))
		return err
	}
	log.Println(fmt.Sprintf("Wizards's %s updated", status))
	return nil
}


// UpdateWizard should update a single status for single Wizard in magicinventory
func UpdateWizardById(w *http.ResponseWriter, db *sql.DB, query string, args  interface{}) (err error){
	var wz dao.Wizard
	row := db.QueryRow(query, args)

	err = row.Scan(&wz.ID, &wz.FirstName, &wz.LastName, &wz.Age, &wz.Category, &wz.Arrested, &wz.Dead)

	if err != nil {
		(*w).WriteHeader(http.StatusNotFound)
		return err
	}

	js, _ := json.Marshal(wz)

	_, err = fmt.Fprintf(*w, string(js))

	if err != nil {
		(*w).WriteHeader(http.StatusUnprocessableEntity)
		log.Println("Cannot update wizard status")
		return err
	}

	log.Printf("Wizards %s 's status have been updated", args)
	return nil
}

// DeleteWizard should update a Wizard in magicinventory
func DeleteWizardById(db *sql.DB, args interface{}) (err error) {
	_, err = db.Exec("DELETE FROM wizards WHERE id = $1;", args)

	if err != nil {
		log.Println("Cannot delete wizard")
		return err
	}

	log.Printf("Wizards %s have been obliviated", args)
	return nil
}

// populateMagicInventory function should create random wizards
// and fill the magicinventory with them
func populateMagicInventory(db *sql.DB) error {
	body, _ := internal.GetRandomNames(10)
	wizards, err := internal.GenerateWizards(body)

	for _, w := range wizards {
		err = CreateWizard(w, db)

		if err != nil {
			return err
		}
	}

	log.Println("Wizards table populated")

	return err
}

// InitMagicInventory function sets up the magicinventory db
// TODO: use `gORM` rather than `pq` ?
// TODO: add an event listener ? https://godoc.org/github.com/lib/pq/example/listen
func InitMagicInventory(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	initQuery :=
		`create table if not exists wizards (
            id         uuid        not null primary key,
            first_name varchar(50) not null,
            last_name  varchar(50) not null,
            age        float       not null,
        	category   varchar(50) not null,
            arrested   boolean     not null,
            dead       boolean     not null
         ); alter table wizards owner to magic;`

	_, err = db.Query(initQuery)

	if err != nil {
		panic(err)
	}

	log.Println("Wizards table created")

	err = populateMagicInventory(db)

	return db, err
}
