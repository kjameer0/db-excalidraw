package main

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_getDrawingByName(t *testing.T) {
	// godotenv.Load(".env.local")
	type fields struct {
		logger *slog.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Valid request with existing drawing",
			fields: fields{
				logger: slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil)),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/drawing?name=example", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				logger: tt.fields.logger,
			}

			// Call the handler
			app.getDrawingByName(tt.args.w, tt.args.r)

			// Verify response
			recorder := tt.args.w.(*httptest.ResponseRecorder)
			if recorder.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
			}
		})
	}
}