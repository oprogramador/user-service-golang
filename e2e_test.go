package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingReadingAddDeleting(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Alan","active":true}`))

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var user User
	assert.Nil(t, err)
	err = json.Unmarshal(bodyBytes, &user)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, "Alan")
	assert.Equal(t, user.Active, true)
	assert.Equal(t, len(user.UserID), 24)

	resp, err = http.Get(ts.URL + "/users")

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	var users []User
	assert.Nil(t, err)
	err = json.Unmarshal(bodyBytes, &users)
	assert.Nil(t, err)
	assert.Contains(t, users, user)

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/user/"+user.UserID, nil)
	assert.Nil(t, err)
	client := &http.Client{}
	resp, err = client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.StatusCode)
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(bodyBytes))

	resp, err = http.Get(ts.URL + "/users")

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(bodyBytes, &users)
	assert.Nil(t, err)
	assert.NotContains(t, users, user)
}

// Thanks to that, DELETE is idempotent
func TestDeletingNonExistentUserWithValidId(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	nonExistentId := "5f31645773bb1c7661d151ba"
	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/user/"+nonExistentId, nil)
	assert.Nil(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(bodyBytes))
}

func TestDeletingWithInvalidId(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/user/invalid", nil)
	assert.Nil(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "invalid id", string(bodyBytes))
}

func TestCreatingUserWithInvalidName(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":123,"active":true}`))

	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "User.name should be of type string", string(bodyBytes))
}

func TestCreatingUserWithInvalidActive(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Alan","active":123}`))

	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "User.active should be of type bool", string(bodyBytes))
}
