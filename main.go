// This application is a demo project that includes simple CRUD features. It is
// used as a demo of gnomock (https://github.com/orlangure/gnomock) abilities,
// and maybe some other things. Check out fedorov.dev for more info.
package main

import (
	"log"
	"net/http"

	"github.com/orlangure/crud-demo/handlers"
	"github.com/orlangure/crud-demo/models"
)

func main() {
	conn := "crud:$tr0ngPasS@/thingdb"

	db, err := models.Connect(conn)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/thing", handlers.CreateThingHandler(db))
	http.HandleFunc("/thing/name", handlers.GetThingByNameHandler(db))
	http.HandleFunc("/thing/id", handlers.GetThingByIDHandler(db))

	log.Fatal(http.ListenAndServe(":8042", nil))
}
