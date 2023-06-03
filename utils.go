package router

import "net/http"

func Vars(r *http.Request) map[string]string {
	ctx := r.Context()
	value := ctx.Value(pathVariablesContextKey).(map[string]string)

	if value == nil {
		return map[string]string{}
	}
	return value
}
