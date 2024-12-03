package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


// Every test should have a proper status code
func Test_application_getDrawingByName(t *testing.T) {
	type fields struct {
		logger    *slog.Logger
		dataSaver dataSaver
	}
	type args struct {
		r *http.Request
	}
	dataSaver := testReader{dataPath: "../../test-drawings"}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Valid drawing retrieval",
			fields: fields{
				logger:    slog.Default(),
				dataSaver: &dataSaver,
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/drawing?name=lorem", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				logger:    tt.fields.logger,
				dataSaver: tt.fields.dataSaver,
			}
			w := httptest.NewRecorder()
			app.getDrawingByName(w, tt.args.r)
			// TODO: change this so that each test has an expected status code instead of hard coding it
			assert.Equal(t, w.Result().StatusCode, http.StatusOK)
		})
	}
}
