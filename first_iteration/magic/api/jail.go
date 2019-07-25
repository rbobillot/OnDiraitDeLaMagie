package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/rbobillo/OnDiraitDeLaMagie/first_iteration/magic/magicinventory"
	"log"
	"net/http"
	"fmt"
)

// UpdateWizardsJail put one wizard in jail by updating
// his jail status.
func UpdateWizardsJail(w *http.ResponseWriter, r *http.Request, db *sql.DB) (err error){
	id := mux.Vars(r)["id"]

	log.Printf("/wizards/{%s}/jail", id)

	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	query := fmt.Sprintf("UPDATE wizards SET arrested = %t WHERE id = $1", true)
	err = magicinventory.UpdateWizardById(db, query, id)

	if err != nil {
		(*w).WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("error: cannot arrest wizards %s", id)
		return err
	}

	return nil
}