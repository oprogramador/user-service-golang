package main

/*
 * these tests are independent on database engine
 * so after replacing it into another one,
 * they will be still passing
 */

import (
	"bytes"
	"encoding/json"
	"github.com/oprogramador/user-service-golang/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlingInvalidActive(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users?active=123")

	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "'active' param in the query string should be of type bool", string(bodyBytes))
}

func TestFilteringByActive(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	_, _ = http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Alan","active":true,"user_id":"3d90f0a3-7109-467a-a0d3-23cfb6fab793"}`))
	_, _ = http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Bob","active":false,"user_id":"cfe765b6-366e-4505-975c-a07dac76b2c3"}`))
	_, _ = http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Cindy","active":true,"user_id":"d87836ff-f360-4b4a-8511-9edaabcda80b"}`))
	_, _ = http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Dave","active":false,"user_id":"7cc4c5f9-c3e2-4778-b775-4ef94d2fc0a0"}`))

	resp, err := http.Get(ts.URL + "/users?active=true")

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var users []models.User
	assert.Nil(t, err)
	err = json.Unmarshal(bodyBytes, &users)
	assert.Nil(t, err)
	assert.Contains(t, users, models.User{Name: "Alan", Active: true, UserID: "3d90f0a3-7109-467a-a0d3-23cfb6fab793"})
	assert.Contains(t, users, models.User{Name: "Cindy", Active: true, UserID: "d87836ff-f360-4b4a-8511-9edaabcda80b"})
	assert.NotContains(t, users, models.User{Name: "Bob", Active: false, UserID: "cfe765b6-366e-4505-975c-a07dac76b2c3"})
	assert.NotContains(t, users, models.User{Name: "Dave", Active: false, UserID: "7cc4c5f9-c3e2-4778-b775-4ef94d2fc0a0"})
}

func TestAddingReadingAddDeleting(t *testing.T) {
	server, _, _, _ := setupServer()
	ts := httptest.NewServer(server)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Alan","active":true}`))

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var user models.User
	assert.Nil(t, err)
	err = json.Unmarshal(bodyBytes, &user)
	assert.Nil(t, err)
	assert.Equal(t, "Alan", user.Name)
	assert.Equal(t, true, user.Active)
	assert.Equal(t, 36, len(user.UserID))

	resp, err = http.Get(ts.URL + "/users")

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	var users []models.User
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

	nonExistentId := "b9cbf5db-81b7-4261-92bd-65f04307b553"
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
