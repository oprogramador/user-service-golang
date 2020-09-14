package main

/*
 * these tests are independent on database engine
 * so after replacing it into another one,
 * they will be still passing
 */

import (
	"bytes"
	"encoding/json"
	. "github.com/franela/goblin"
	"github.com/google/uuid"
	"github.com/oprogramador/user-service-golang/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	g := Goblin(t)
	g.Describe("routes", func() {
		g.It("handles invalid active", func() {
			// GIVEN
			server, _, _, _ := setupServer()
			ts := httptest.NewServer(server)
			defer ts.Close()

			// WHEN
			resp, err := http.Get(ts.URL + "/users?active=123")

			// THEN
			assert.Nil(t, err)
			assert.Equal(t, 400, resp.StatusCode)
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, "'active' param in the query string should be of type bool", string(bodyBytes))
		})

		g.It("filters by active", func() {
			// GIVEN
			server, _, _, _ := setupServer()
			ts := httptest.NewServer(server)
			defer ts.Close()
			id1, err := uuid.NewRandom()
			assert.Nil(t, err)
			id2, err := uuid.NewRandom()
			assert.Nil(t, err)
			id3, err := uuid.NewRandom()
			assert.Nil(t, err)
			id4, err := uuid.NewRandom()
			assert.Nil(t, err)

			_, _ = http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Alan","active":true,"user_id":"`+id1.String()+`"}`))
			_, _ = http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Bob","active":false,"user_id":"`+id2.String()+`"}`))
			_, _ = http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Cindy","active":true,"user_id":"`+id3.String()+`"}`))
			_, _ = http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Dave","active":false,"user_id":"`+id4.String()+`"}`))

			// WHEN
			resp, err := http.Get(ts.URL + "/users?active=true")

			// THEN
			assert.Nil(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			var users []models.User
			assert.Nil(t, err)
			err = json.Unmarshal(bodyBytes, &users)
			assert.Nil(t, err)
			assert.Contains(t, users, models.User{Name: "Alan", Active: true, UserID: id1.String()})
			assert.Contains(t, users, models.User{Name: "Cindy", Active: true, UserID: id3.String()})
			assert.NotContains(t, users, models.User{Name: "Bob", Active: false, UserID: id2.String()})
			assert.NotContains(t, users, models.User{Name: "Dave", Active: false, UserID: id4.String()})
		})

		g.It("adds, reads and deletes", func() {
			// GIVEN
			server, _, _, _ := setupServer()
			ts := httptest.NewServer(server)
			defer ts.Close()

			// WHEN
			resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Alan","active":true}`))

			// THEN
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

			// WHEN
			resp, err = http.Get(ts.URL + "/users")

			// THEN
			assert.Nil(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			var users []models.User
			assert.Nil(t, err)
			err = json.Unmarshal(bodyBytes, &users)
			assert.Nil(t, err)
			assert.Contains(t, users, user)

			// WHEN
			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/user/"+user.UserID, nil)
			assert.Nil(t, err)
			client := &http.Client{}
			resp, err = client.Do(req)

			// THEN
			assert.Nil(t, err)
			assert.Equal(t, 204, resp.StatusCode)
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, "", string(bodyBytes))

			// WHEN
			resp, err = http.Get(ts.URL + "/users")

			// THEN
			assert.Nil(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, []string{"application/json; charset=utf-8"}, resp.Header["Content-Type"])
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(bodyBytes, &users)
			assert.Nil(t, err)
			assert.NotContains(t, users, user)
		})

		// Thanks to that, DELETE is idempotent
		g.It("deletes non-existent user", func() {
			// GIVEN
			server, _, _, _ := setupServer()
			ts := httptest.NewServer(server)
			defer ts.Close()
			nonExistentId := "b9cbf5db-81b7-4261-92bd-65f04307b553"

			// WHEN
			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/user/"+nonExistentId, nil)
			assert.Nil(t, err)
			client := &http.Client{}
			resp, err := client.Do(req)

			// THEN
			assert.Nil(t, err)
			assert.Equal(t, 204, resp.StatusCode)
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, "", string(bodyBytes))
		})

		g.It("creates user with invalid name", func() {
			// GIVEN
			server, _, _, _ := setupServer()
			ts := httptest.NewServer(server)
			defer ts.Close()

			// WHEN
			resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":123,"active":true}`))

			// THEN
			assert.Nil(t, err)
			assert.Equal(t, 400, resp.StatusCode)
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, "User.name should be of type string", string(bodyBytes))
		})

		g.It("forbids creating user with invalid active", func() {
			// GIVEN
			server, _, _, _ := setupServer()
			ts := httptest.NewServer(server)
			defer ts.Close()

			// WHEN
			resp, err := http.Post(ts.URL+"/user", "application/json", bytes.NewBufferString(`{"name":"Alan","active":123}`))

			// THEN
			assert.Nil(t, err)
			assert.Equal(t, 400, resp.StatusCode)
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, "User.active should be of type bool", string(bodyBytes))
		})

	})
}
