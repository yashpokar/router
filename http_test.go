package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stringToMethod(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want method
	}{
		{
			name: "returns GET",
			args: args{
				name: "GET",
			},
			want: GET,
		},
		{
			name: "returns GET",
			args: args{
				name: "GET",
			},
			want: GET,
		},
		{
			name: "returns POST",
			args: args{
				name: "POST",
			},
			want: POST,
		},
		{
			name: "returns PUT",
			args: args{
				name: "PUT",
			},
			want: PUT,
		},
		{
			name: "returns PATCH",
			args: args{
				name: "PATCH",
			},
			want: PATCH,
		},
		{
			name: "returns DELETE",
			args: args{
				name: "DELETE",
			},
			want: DELETE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, stringToMethod(tt.args.name), "stringToMethod(%v)", tt.args.name)
		})
	}

	t.Run("will panic for unsupported method", func(t *testing.T) {
		assert.Panics(t, func() {
			stringToMethod("HEAD")
		}, "method HEAD is not supported")
	})
}
