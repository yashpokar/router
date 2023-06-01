package router

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yashpokar/router/tests"
)

func TestGroup(t *testing.T) {
	r := New()
	api := r.Group("/api")

	server := http.Server{Handler: r, Addr: ":8989"}
	go func() {
		err := server.ListenAndServe()
		assert.NoError(t, err)
	}()

	t.Run("registers GET method", func(t *testing.T) {
		api.GET("/", tests.Handler)
	})

	t.Run("registers POST method", func(t *testing.T) {
		api.POST("/", tests.Handler)
	})

	t.Run("registers PUT method", func(t *testing.T) {
		api.PUT("/", tests.Handler)
	})

	t.Run("registers PATCH method", func(t *testing.T) {
		api.PATCH("/", tests.Handler)
	})

	t.Run("registers DELETE method", func(t *testing.T) {
		api.DELETE("/", tests.Handler)
	})

	t.Run("registers GET method on nested group", func(t *testing.T) {
		v1 := api.Group("/v1")
		v1.GET("/", tests.Handler)
	})

	t.Run("resolves the nested GET method", func(t *testing.T) {
		v1 := api.Group("/v1")
		usersGroup := v1.Group("/users")
		usersGroup.GET("/", tests.Handler)

		response, err := http.Get("http://localhost:8989/api/v1/users/")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("returns 404 status when the route does not exists for the uri", func(t *testing.T) {
		response, err := http.Get("http://localhost:8989/api/v1/unknown/route")
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			assert.NoError(t, err)
		}(response.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}
