package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/orlangure/crud-demo/handlers"
	"github.com/orlangure/crud-demo/models"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/mysql"
	"github.com/stretchr/testify/require"
)

var db *models.DB

func TestCRUD(t *testing.T) {
	data := url.Values{}
	data.Add("name", "thing")
	data.Add("comment", "this is a thing")
	body := bytes.NewBufferString(data.Encode())

	r, w := httptest.NewRequest(http.MethodPost, "/thing", body), httptest.NewRecorder()
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	createHandler := handlers.CreateThingHandler(db)
	createHandler(w, r)

	res := w.Result()

	resBody, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equalf(t, http.StatusCreated, res.StatusCode, string(resBody))

	r, w = httptest.NewRequest(http.MethodGet, "/thing/name?name=thing", nil), httptest.NewRecorder()

	byNameHandler := handlers.GetThingByNameHandler(db)
	byNameHandler(w, r)

	res = w.Result()

	resBody, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equalf(t, http.StatusOK, res.StatusCode, string(resBody))

	thing := &models.Thing{}
	require.NoError(t, json.Unmarshal(resBody, &thing))
	require.Equal(t, "thing", thing.Name)
	require.Equal(t, "this is a thing", thing.Comment)

	url := fmt.Sprintf("/thing/id?id=%d", thing.ID)
	r, w = httptest.NewRequest(http.MethodGet, url, nil), httptest.NewRecorder()

	byIDHandler := handlers.GetThingByIDHandler(db)
	byIDHandler(w, r)

	res = w.Result()

	resBody, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equalf(t, http.StatusOK, res.StatusCode, string(resBody))

	thing = &models.Thing{}
	require.NoError(t, json.Unmarshal(resBody, &thing))
	require.Equal(t, "thing", thing.Name)
	require.Equal(t, "this is a thing", thing.Comment)
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) (code int) {
	p := mysql.Preset(
		mysql.WithDatabase("thingdb"),
		mysql.WithUser("crud", "password"),
		mysql.WithQueriesFile("./schema/schema.sql"),
	)

	c, err := gnomock.Start(p)
	if err != nil {
		log.Println(err)
		return 1
	}

	defer func() {
		err := gnomock.Stop(c)
		if err != nil {
			log.Println(err)

			if code == 0 {
				code = 3
			}
		}
	}()

	conn := fmt.Sprintf("crud:password@tcp(%s)/thingdb", c.DefaultAddress())

	db, err = models.Connect(conn)
	if err != nil {
		log.Println(err)
		return 2
	}

	return m.Run()
}
