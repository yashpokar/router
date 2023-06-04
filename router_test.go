package router

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yashpokar/router/tests"
)

func ProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := Vars(r)
	productID, ok := vars["product_id"]
	if !ok {
		_, err := w.Write([]byte("'product_id' not found"))
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err := w.Write([]byte(productID))
	if err != nil {
		log.Fatal(err)
	}
}

func TestRouter(t *testing.T) {
	r := New().(*router)
	server := &http.Server{Handler: r, Addr: ":8787"}
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
		r.handleError(tests.NewMockResponseWriter(), &http.Request{}, errors.New("internal server error"))
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

	t.Run("returns error when route not found", func(t *testing.T) {
		route, err := r.resolve("GET", "unknown/path")

		assert.EqualError(t, err, "route 'unknown/path' not found")
		assert.Nil(t, route)
	})

	t.Run("returns the product id using path variable", func(t *testing.T) {
		r.GET("/products/:product_id", ProductDetailsHandler)

		response, err := http.Get("http://localhost:8787/products/8298ae5f-9cbf-4751-a411-560419b0b5d7")
		assert.NoError(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, []byte("8298ae5f-9cbf-4751-a411-560419b0b5d7"), bytes)
		defer response.Body.Close()
	})

	t.Run("can not return the product id when it is not there", func(t *testing.T) {
		r.GET("/products", ProductDetailsHandler)

		response, err := http.Get("http://localhost:8787/products")
		assert.NoError(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.Equal(t, []byte("'product_id' not found"), bytes)
		defer response.Body.Close()
	})

	t.Run("recovers from the panic", func(t *testing.T) {
		r.GET("/panic-maker", tests.PanicMaker)
		r.OnPanic(func(w http.ResponseWriter, _ *http.Request, r any) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(r.(string)))
		})

		response, err := http.Get("http://localhost:8787/panic-maker")
		assert.NoError(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		assert.Equal(t, []byte("to test panic recovery"), bytes)
		defer response.Body.Close()
	})

	t.Run("handles the not found route", func(t *testing.T) {
		r.GET("/panic-maker", tests.PanicMaker)
		r.OnRouteNotFound(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fmt.Sprintf("route '%s' not found.", request.RequestURI)))
		})

		response, err := http.Get("http://localhost:8787/route/not/found")
		assert.NoError(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.Equal(t, []byte("route '/route/not/found' not found."), bytes)
		defer response.Body.Close()
	})

	t.Run("handles the default error when panic handler is provided", func(t *testing.T) {
		r.OnPanic(func(_ http.ResponseWriter, _ *http.Request, err any) {
			e, ok := err.(error)
			assert.True(t, ok)
			assert.EqualError(t, e, tests.ErrUnknown.Error())
		})

		r.handleError(tests.NewMockResponseWriter(), &http.Request{}, tests.ErrUnknown)
	})
}
