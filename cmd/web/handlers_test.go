package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// needs for tests: i want to be able to have tests to test generic http stuff but there are special tests maybe
func Test_application_getDrawingByName(t *testing.T) {
	type fields struct {
		logger    *slog.Logger
		dataSaver dataSaver
	}
	type args struct {
		r            *http.Request
		expectedCode int
	}
	dataSaver := testReader{dataPath: "../../test-drawings"}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "DrawingFound",
			fields: fields{
				logger:    logger,
				dataSaver: &dataSaver,
			},
			args: args{
				r:            httptest.NewRequest(http.MethodGet, "/drawing?name=lorem", nil),
				expectedCode: http.StatusOK,
			},
		},
		{
			name: "DrawingNotFound",
			fields: fields{
				logger:    logger,
				dataSaver: &dataSaver,
			},
			args: args{
				r:            httptest.NewRequest(http.MethodGet, "/drawing?name=lore", nil),
				expectedCode: http.StatusNotFound,
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
			assert.Equal(t, tt.args.expectedCode, w.Result().StatusCode)
		})
	}
}
