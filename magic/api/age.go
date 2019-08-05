package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rbobillo/OnDiraitDeLaMagie/magic/dao"
	"github.com/rbobillo/OnDiraitDeLaMagie/magic/internal"
	"github.com/rbobillo/OnDiraitDeLaMagie/magic/magicinventory"
	"net/http"
)

// AgeWizards function request the Magic Inventory to update every wizard age by increment it n times
func AgeWizards(w *http.ResponseWriter, r *http.Request, db *sql.DB) (err error) {
	var wizard dao.Wizard

	internal.Debug("/wizards/age")

	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&wizard)

	if err != nil {
		(*w).WriteHeader(http.StatusMethodNotAllowed)
		internal.Warn("cannot convert Body to JSON")
		return err
	}

	query := "UPDATE wizards SET age = age + $1;"
	err = magicinventory.UpdateWizards(db, query, wizard.Age)

	if err != nil {
		(*w).WriteHeader(http.StatusUnprocessableEntity)
		internal.Warn("cannot update wizards's age")
		return err
	}

	(*w).WriteHeader(http.StatusNoContent)
	internal.Info(fmt.Sprintf("every wizards got older (+%f years)", wizard.Age))

	return nil
}