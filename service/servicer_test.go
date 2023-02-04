package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/claco/claco.pinochle/pb"
	"github.com/claco/claco.pinochle/service"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	name     string
	user     string
	password string
	instance testcontainers.Container
	host     string
	port     int
}

func (db *TestDatabase) Host() string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ip, _ := db.instance.Host(ctx)

	return ip
}

func (db *TestDatabase) Port() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	port, _ := db.instance.MappedPort(ctx, "5432/tcp")

	return port.Int()
}

// set to servicer code defaults
var database = &TestDatabase{
	name:     "pinochle",
	user:     "pinochle",
	password: "pinochle",
	host:     "localhost",
	port:     5432,
}

func TestMain(m *testing.M) {
	os.Chdir("../") //nolint

	if _, integration := os.LookupEnv("INTEGRATION"); integration == true {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		// set to postgres defaults prior to migrations
		database.name = "postgres"
		database.user = "postgres"
		database.password = "postgres"

		request := testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			SkipReaper:   true,
			AutoRemove:   true,
			Env: map[string]string{
				"POSTGRES_USER":     database.user,
				"POSTGRES_PASSWORD": database.password,
				"POSTGRES_DB":       database.name,
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})
		defer container.Terminate(ctx) //nolint

		if err != nil {
			os.Exit(1)
		}

		database.instance = container
		database.host = database.Host()
		database.port = database.Port()

		os.Setenv("DATABASE_HOST", database.host)
		os.Setenv("DATABASE_PORT", fmt.Sprint((database.port)))
		os.Setenv("DATABASE_NAME", database.name)
		os.Setenv("DATABASE_USER", database.user)
		os.Setenv("DATABASE_PASSWORD", database.password)

		m.Run()

		if container.IsRunning() {
			container.Terminate(ctx) //nolint
		}
	} else {
		m.Run()
	}
}

func TestPinichleService_DatabaseConnectionString(t *testing.T) {
	type args struct {
		showPassword bool
	}

	svc := service.NewPinochleService()

	tests := []struct {
		name string
		args args
		want string
	}{
		{"hides password by default", args{false}, fmt.Sprintf("postgres://%s:*****@localhost:%d/%s?application_name=pinochle", database.user, database.port, database.name)},
		{"shows password when requested", args{true}, fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?application_name=pinochle", database.user, database.password, database.port, database.name)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.DatabaseConnectionString(tt.args.showPassword); got != tt.want {
				t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("reads DATABASE_USER from env", func(t *testing.T) {
		user := "test-user"
		want := fmt.Sprintf("postgres://%s:*****@localhost:%d/%s?application_name=pinochle", user, database.port, database.name)
		t.Setenv("DATABASE_USER", user)

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("returns default database user", func(t *testing.T) {
		user := "pinochle"
		want := fmt.Sprintf("postgres://%s:*****@localhost:%d/%s?application_name=pinochle", user, database.port, database.name)
		t.Setenv("DATABASE_USER", "")

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("reads DATABASE_PASSWORD from env", func(t *testing.T) {
		password := "test-password"
		want := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?application_name=pinochle", database.user, password, database.port, database.name)
		t.Setenv("DATABASE_PASSWORD", password)

		got := svc.DatabaseConnectionString(true)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("returns default database password", func(t *testing.T) {
		password := "pinochle"
		want := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?application_name=pinochle", database.user, password, database.port, database.name)
		t.Setenv("DATABASE_PASSWORD", "")

		got := svc.DatabaseConnectionString(true)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("reads DATABASE_HOST from env", func(t *testing.T) {
		host := "test-host"
		want := fmt.Sprintf("postgres://%s:*****@%s:%d/%s?application_name=pinochle", database.user, host, database.port, database.name)
		t.Setenv("DATABASE_HOST", host)

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("returns default database host", func(t *testing.T) {
		host := "localhost"
		want := fmt.Sprintf("postgres://%s:*****@%s:%d/%s?application_name=pinochle", database.user, host, database.port, database.name)
		t.Setenv("DATABASE_HOST", "")

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("reads DATABASE_PORT from env", func(t *testing.T) {
		port := "1234"
		want := fmt.Sprintf("postgres://%s:*****@%s:%s/%s?application_name=pinochle", database.user, database.host, port, database.name)
		t.Setenv("DATABASE_PORT", port)

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("returns default database port", func(t *testing.T) {
		port := "5432"
		want := fmt.Sprintf("postgres://%s:*****@%s:%s/%s?application_name=pinochle", database.user, database.host, port, database.name)
		t.Setenv("DATABASE_PORT", "")

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("reads DATABASE_NAME from env", func(t *testing.T) {
		name := "test-name"
		want := fmt.Sprintf("postgres://%s:*****@%s:%v/%s?application_name=pinochle", database.user, database.host, database.port, name)
		t.Setenv("DATABASE_NAME", name)

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("returns default database name", func(t *testing.T) {
		name := "pinochle"
		want := fmt.Sprintf("postgres://%s:*****@%s:%v/%s?application_name=pinochle", database.user, database.host, database.port, name)
		t.Setenv("DATABASE_NAME", "")

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("reads DATABASE_APPLICATION_NAME from env", func(t *testing.T) {
		name := "test-name"
		want := fmt.Sprintf("postgres://%s:*****@%s:%v/%s?application_name=%s", database.user, database.host, database.port, database.name, name)
		t.Setenv("DATABASE_APPLICATION_NAME", name)

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("returns default database application name", func(t *testing.T) {
		name := "pinochle"
		want := fmt.Sprintf("postgres://%s:*****@%s:%v/%s?application_name=%s", database.user, database.host, database.port, database.name, name)
		t.Setenv("DATABASE_APPLICATION_NAME", "")

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})

	t.Run("appends DATABASE_SCHEMA from env to search path", func(t *testing.T) {
		schema := "test-schema"
		want := fmt.Sprintf("postgres://%s:*****@%s:%v/%s?application_name=pinochle&search_path=%s,public", database.user, database.host, database.port, database.name, schema)
		t.Setenv("DATABASE_SCHEMA", schema)

		got := svc.DatabaseConnectionString(false)

		if got != want {
			t.Errorf("PinichleService.DatabaseConnectionString() = %v, want %v", got, want)
		}
	})
}

func TestPinichleService_GetKey(t *testing.T) {
	type args struct {
		msg service.PinochleRecordLookup
	}

	svc := service.NewPinochleService()
	tests := []struct {
		name string
		args args
		want string
	}{
		{"returns slug when id is empty", args{msg: &pb.Game{Slug: "my-slug"}}, "my-slug"},
		{"returns id when slug is empty", args{msg: &pb.Game{Id: "00000000-0000-0000-0000-000000000000"}}, "00000000-0000-0000-0000-000000000000"},
		{"returns slug when both are set", args{msg: &pb.Game{Id: "00000000-0000-0000-0000-000000000000", Slug: "my-slug"}}, "my-slug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.GetKey(tt.args.msg); got != tt.want {
				t.Errorf("PinichleService.GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPinichleService_Slugify(t *testing.T) {
	type args struct {
		name string
	}

	svc := service.NewPinochleService()
	tests := []struct {
		name string
		args args
		want string
	}{
		{"converts spaces to deshes", args{name: "my slug"}, "my-slug"},
		{"splits pascal case words", args{name: "MySlug"}, "my-slug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.Slugify(tt.args.name); got != tt.want {
				t.Errorf("PinichleService.Slugify() = %v, want %v", got, tt.want)
			}
		})
	}

	fmt.Println(time.Now())
}

func TestPinichleService_ConnectDatabase(t *testing.T) {
	svc := service.NewPinochleService()

	if _, integration := os.LookupEnv("INTEGRATION"); integration == false {
		t.Skip("set env[INTEGRATION] to enable")
		t.SkipNow()
	}

	t.Run("can connect to the database", func(t *testing.T) {
		if err := svc.ConnectDatabase(); err != nil {
			t.Error(err)
		}
	})

	t.Run("handles database connection errors", func(t *testing.T) {
		t.Setenv("DATABASE_USER", "non-existent")
		t.Setenv("DATABASE_PASSWORD", "invalid-password")
		svc.CloseDatabase()

		if err := svc.ConnectDatabase(); err == nil {
			t.Log(err)
			t.Error("PinichleService.ConnectDatabase() accidentally succeeded with invalid user/password")
		}
	})
}

func TestPinichleService_InitializeDatabase(t *testing.T) {
	svc := service.NewPinochleService()

	if _, integration := os.LookupEnv("INTEGRATION"); integration == false {
		t.Skip("set env[INTEGRATION] to enable")
		t.SkipNow()
	}

	if err := svc.ConnectDatabase(); err != nil {
		t.Fatal(err)
	}

	t.Run("can initialize the database", func(t *testing.T) {
		if err := svc.InitializeDatabase(); err != nil {
			t.Error(err)
		}
	})

	t.Run("returns error for missing migrations", func(t *testing.T) {
		cwd, _ := os.Getwd() //nolint
		defer os.Chdir(cwd)  //nolint

		os.Chdir(t.TempDir()) //nolint

		if err := svc.InitializeDatabase(); err == nil {
			t.Error("PinichleService.InitializeDatabase() accidentally succeeded with missing migrations")
		}
	})
}
