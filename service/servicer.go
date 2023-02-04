package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/claco/claco.pinochle/pb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type PinichleService struct {
	pb.UnimplementedPinochleServiceServer
	mu *sync.Mutex
	db *pgx.Conn
}

type PinochleRecordLookup interface {
	GetId() string
	GetSlug() string
}

func NewPinochleService() *PinichleService {
	service := &PinichleService{}
	service.mu = &sync.Mutex{}

	return service
}

func (svc *PinichleService) ConnectDatabase() error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if svc.db == nil {
		conn, err := pgx.Connect(context.Background(), svc.DatabaseConnectionString(true))

		maskedConnectionString := svc.DatabaseConnectionString(false)
		log.WithField("database", maskedConnectionString).Logger.Infof("connecting to database: %s", maskedConnectionString)

		if err != nil {
			return err
		} else {
			svc.db = conn
		}
	}

	return nil
}

func (svc *PinichleService) CloseDatabase() {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if svc.db != nil {
		svc.db.Close(context.Background())
		svc.db = nil
	}
}

func (svc *PinichleService) GetKey(msg PinochleRecordLookup) string {
	id, slug := msg.GetId(), msg.GetSlug()

	if slug != "" {
		return slug
	} else {
		return id
	}
}

func (svc *PinichleService) InitializeDatabase() error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	database_schema, database_schema_set := os.LookupEnv("DATABASE_SCHEMA")
	database_url := svc.DatabaseConnectionString(true)
	database_url = strings.Replace(database_url, "postgres://", "pgx://", 1)

	if !database_schema_set {
		current_user := svc.db.Config().User
		user_schema := current_user
		fields := log.Fields{
			"user":   current_user,
			"schema": user_schema,
		}

		log.WithFields(fields).Infof("creating user schema: %s", user_schema)
		_, err := svc.db.Exec(context.Background(), fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS AUTHORIZATION %s", current_user))
		if err != nil {
			return err
		}
	} else {
		log.WithField("schema", database_schema).Infof("using existing schema: %s", database_schema)
	}

	maskedConnectionString := svc.DatabaseConnectionString(false)
	log.WithField("database", maskedConnectionString).Debugf("applying database migrations: %s", maskedConnectionString)

	migrations, err := migrate.New("file://db/migrate", database_url)

	if err != nil {
		return err
	} else {
		err := migrations.Up()
		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Debug("no migration changes were applied")
			} else {
				return err
			}
		} else {
			log.Info("database migrations applied")
		}
	}

	return nil
}

func (svc *PinichleService) DatabaseConnectionString(showPassword bool) string {
	database_schema, database_schema_set := os.LookupEnv("DATABASE_SCHEMA")

	database_application_name, ok := os.LookupEnv("DATABASE_APPLICATION_NAME")
	if !ok || database_application_name == "" {
		database_application_name = "pinochle"
	}
	database_host, ok := os.LookupEnv("DATABASE_HOST")
	if !ok || database_host == "" {
		database_host = "localhost"
	}
	database_port, ok := os.LookupEnv("DATABASE_PORT")
	if !ok || database_port == "" {
		database_port = "5432"
	}
	database_name, ok := os.LookupEnv("DATABASE_NAME")
	if !ok || database_name == "" {
		database_name = "pinochle"
	}
	database_user, ok := os.LookupEnv("DATABASE_USER")
	if !ok || database_user == "" {
		database_user = "pinochle"
	}
	database_password, ok := os.LookupEnv("DATABASE_PASSWORD")
	if !ok || database_password == "" {
		database_password = "pinochle"
	}
	if !showPassword {
		database_password = "*****"
	}

	database_url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?application_name=%s", database_user, database_password, database_host, database_port, database_name, database_application_name)
	if database_schema_set && database_schema != "" {
		database_url = fmt.Sprintf("%s&search_path=%s,public", database_url, database_schema)
	}

	return database_url
}

func (svc *PinichleService) Slugify(name string) string {
	words := svc.SplitPascalCase(name)

	return slug.Make(words)
}

func (svc *PinichleService) SplitPascalCase(name string) string {
	var buffer bytes.Buffer
	var priorLower bool = false

	for _, char := range name {
		if priorLower && unicode.IsUpper(char) {
			buffer.WriteByte(' ')
		}
		buffer.WriteRune(char)
		priorLower = unicode.IsLower(char)
	}

	return buffer.String()
}
