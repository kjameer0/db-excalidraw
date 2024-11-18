package internal

import (
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestUuidToHash(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("env failed to load")
	}
	type args struct {
		uuid uuid.UUID
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Test valid UUID 1",
			args: args{
				uuid: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			want:  "someExpectedHash1",
			want1: "anotherExpectedHash1",
		},
		{
			name: "Test valid UUID 2",
			args: args{
				uuid: uuid.MustParse("9a4e2a70-59c3-11ec-bf63-0242ac130002"),
			},
			want:  "someExpectedHash2",
			want1: "anotherExpectedHash2",
		},
		{
			name: "Test valid UUID with specific edge case",
			args: args{
				uuid: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
			want:  "edgeCaseHash1",
			want1: "edgeCaseHash2",
		},
		{
			name: "Test valid UUID with all F's",
			args: args{
				uuid: uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
			},
			want:  "allFHash1",
			want1: "allFHash2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UuidToHash(tt.args.uuid, []byte(os.Getenv("AES_SECRET")))
			if got != tt.want {
				t.Errorf("UuidToHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}
