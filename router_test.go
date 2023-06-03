package router

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yashpokar/router/tests"
)

func ProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := Vars(r)
	productID, ok := vars["product_id"]
	if !ok {
		w.Write([]byte("'product_id' not found"))
	}

	w.Write([]byte(productID))
}

func TestRouter(t *testing.T) {
	r := New().(*router)
	server := http.Server{Handler: r, Addr: ":8787"}
	go func() {
		err := server.ListenAndServe()
		assert.NoError(t, err)
	}()

	t.Run("registers GET method", func(t *testing.T) {
		r.GET("/", tests.Handler)
		route, err := r.resolve("GET", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("registers POST method", func(t *testing.T) {
		r.POST("/", tests.Handler)
		route, err := r.resolve("POST", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("registers PUT method", func(t *testing.T) {
		r.PUT("/", tests.Handler)
		route, err := r.resolve("PUT", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("registers PATCH method", func(t *testing.T) {
		r.PATCH("/", tests.Handler)
		route, err := r.resolve("PATCH", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("registers DELETE method", func(t *testing.T) {
		r.DELETE("/", tests.Handler)
		route, err := r.resolve("DELETE", "")
		assert.NoError(t, err)
		assert.Equal(t, route.path, "")
	})

	t.Run("handles general error", func(t *testing.T) {
		r.handleError(errors.New("internal server error"), tests.NewMockResponseWriter())
	})

	t.Run("resolves the route", func(t *testing.T) {
		r.GET("/path/nested", tests.Handler)
		route, err := r.resolve("GET", "path/nested")

		assert.NoError(t, err)
		assert.Equal(t, "path/nested", route.path)
	})

	t.Run("returns the path variables map", func(t *testing.T) {
		r.GET("/products/:product_id", tests.Handler)
		route, err := r.resolve("GET", "products/8298ae5f-9cbf-4751-a411-560419b0b5d7")

		assert.NoError(t, err)
		assert.Equal(t, route.getPathVariablesMap(), map[string]string{"product_id": "8298ae5f-9cbf-4751-a411-560419b0b5d7"})
	})

	t.Run("returns nil when there are no path variables", func(t *testing.T) {
		r.GET("/products", tests.Handler)
		route, err := r.resolve("GET", "products")

		assert.NoError(t, err)
		assert.Nil(t, route.getPathVariablesMap())
	})

	t.Run("returns the product id using path variable", func(t *testing.T) {
		r.GET("/products/:product_id", ProductDetailsHandler)

		response, err := http.Get("http://localhost:8787/products/8298ae5f-9cbf-4751-a411-560419b0b5d7")
		bytes, err := io.ReadAll(response.Body)
		assert.NoError(t, err)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			assert.NoError(t, err)
		}(response.Body)

		assert.NoError(t, err)
		assert.Equal(t, []byte("8298ae5f-9cbf-4751-a411-560419b0b5d7"), bytes)
	})

	t.Run("can not return the product id when it is not there", func(t *testing.T) {
		r.GET("/products", ProductDetailsHandler)

		response, err := http.Get("http://localhost:8787/products")
		bytes, err := io.ReadAll(response.Body)
		assert.NoError(t, err)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			assert.NoError(t, err)
		}(response.Body)

		assert.NoError(t, err)
		assert.Equal(t, []byte("'product_id' not found"), bytes)
	})
}
