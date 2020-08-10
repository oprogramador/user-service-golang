package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingAndReading(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Alan","active":true}`))

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var user User
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, user.Name, "Alan")
	assert.Equal(t, user.Active, true)
	assert.Equal(t, len(user.UserID), 24)

	resp, err = http.Get(ts.URL + "/users")

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	var users []User
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &users)
	if err != nil {
		log.Fatal(err)
	}
	assert.Contains(t, users, user)
}
