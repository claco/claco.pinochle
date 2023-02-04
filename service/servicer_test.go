package service

import (
	"sync"
	"testing"

	"github.com/claco/claco.pinochle/pb"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func TestPinichleService_Slugify(t *testing.T) {
	type fields struct {
		UnimplementedPinochleServiceServer pb.UnimplementedPinochleServiceServer
		mu                                 *sync.Mutex
		db                                 *pgx.Conn
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{"converts spaces to deshes", fields{}, args{name: "my slug"}, "my-slug"},
		{"splits pascal case words", fields{}, args{name: "MySlug"}, "my-slug"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := &PinichleService{
				UnimplementedPinochleServiceServer: tt.fields.UnimplementedPinochleServiceServer,
				mu:                                 tt.fields.mu,
				db:                                 tt.fields.db,
			}
			if got := svc.Slugify(tt.args.name); got != tt.want {
				t.Errorf("PinichleService.Slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}
