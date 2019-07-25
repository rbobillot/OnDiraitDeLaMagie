package api

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rbobillo/OnDiraitDeLaMagie/first_iteration/magic/magicinventory"
	"log"
	"net/http"
)

// UpdateWizardsAges function request the Magic Inventory
// to update one wizard
// Todo: UpdateWizardDeath and UpdateWizardJail are almost the same function
//
func UpdateWizardsDeath(w *http.ResponseWriter, r *http.Request, db *sql.DB) (err error){

	id := mux.Vars(r)["id"]

	log.Printf("/wizards/{%s}/die", id)

	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")

	query := fmt.Sprintf("UPDATE wizards SET dead = %t WHERE id = $1 RETURNING *;", true)
	err = magicinventory.UpdateWizardById(db, query, id)

	if err != nil {
		(*w).WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("error: cannot kill wizards %s", id)
		return err
	}

	return nil
}
