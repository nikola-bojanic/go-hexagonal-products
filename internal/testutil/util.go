package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/emicklei/go-restful/v3"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/app"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/database"
	"github.com/pkg/errors"
)

func InitTestApp() *app.App {
	app := app.MustInitializeTestApp()

	migration := database.NewMigrationProcess(app.DB, app.Logger)

	err := migration.DropSchema(app.Config.Database.Schema)
	if err != nil {
		panic(errors.Wrap(err, "error resetting DB"))
	}

	// todo: fix migrations: tests can be 2 or 3 subdirs deep - check for migrations folder existence before migrating
	err = migration.Migrate("../../../migrations")
	if err != nil {
		panic(errors.Wrap(err, "error migrating DB"))
	}

	return app
}

// Sends a HTTP request to the set path, using the set method
// If postData is provided, it will be marshalled and sent
func MakeRequest(container *restful.Container, method string, path string, postData interface{}, jwtToken *string) *httptest.ResponseRecorder {
	jsdat, _ := json.Marshal(postData)
	bodyReader := bytes.NewReader(jsdat)
	httpRequest, _ := http.NewRequest(method, path, bodyReader)
	httpRequest.Header.Set("Content-Type", restful.MIME_JSON)
	if jwtToken != nil {
		httpRequest.Header.Set("Authorization", "Bearer "+(*jwtToken))
	}
	responseRec := httptest.NewRecorder()

	// send request
	container.ServeHTTP(responseRec, httpRequest)

	return responseRec
}

// Deletes all records from all tables
func CleanUpTables(db database.DB) {
	db.Exec(context.TODO(), "TRUNCATE TABLE hex_fwk.user")
	db.Exec(context.TODO(), "TRUNCATE TABLE hex_fwk.product")
	db.Exec(context.TODO(), "TRUNCATE TABLE hex_fwk.category")
}
