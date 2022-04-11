// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"finstar-test-task/proto"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Postgres -.
type Postgres struct {
	database *gorm.DB
}

func (p Postgres) GetDatabase() *gorm.DB {
	return p.database
}

type RecorderLogger struct {
	logger.Interface
	Statements []string
}

func (r *RecorderLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()

	if !strings.HasPrefix(sql, "SELECT ") {
		log.Printf("Going to save '%s' as new up migration\n", sql)

		beginStringified := begin.UTC().Format("20060102150405")
		vital := strings.Join(strings.Split(sql, " ")[0:3], "_")
		raw := fmt.Sprintf("migrations/%s_%s.up.sql", beginStringified, vital)
		raw = strings.ToLower(strings.Replace(raw, "\"", "", -1))
		content, err := os.Create(raw)
		if err != nil {
			panic(fmt.Errorf("os create failed: %v\n", err))
		}
		defer content.Close()

		if _, err = content.WriteString(fmt.Sprintf("%s;\n", sql)); err != nil {
			panic(fmt.Errorf("write string failed: %v\n", err))
		}
	}

	r.Statements = append(r.Statements, sql)
}

// New -.
func New(dsn string, logLevel int, autoMigrate, generateSeeds bool) (*Postgres, error) {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.LogLevel(logLevel),
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("Couldn't establish database connection: %s", err.Error())
	}

	log.Printf("Postgres database connection established!")

	if sqlDb, err := db.DB(); err != nil {
		log.Fatal(err)
	} else if sqlDb.Ping() != nil {
		log.Fatal(err)
	}

	if autoMigrate {
		for _, model := range []interface{}{
			proto.UserORM{},
		} {
			recorder := RecorderLogger{}
			recorder.Interface = logger.Default.LogMode(logger.Info)
			session := db.Session(&gorm.Session{
				Logger: &recorder,
			})
			session.AutoMigrate(model)

			if !db.Migrator().HasTable(model) {
				if err := db.AutoMigrate(model); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	if generateSeeds {
		generateUsers(db, 2)
	}

	return &Postgres{database: db}, nil
}

func (pg Postgres) Close() {
	log.Println("Close Postgres DB connection...")
	var err error
	db, err := pg.database.DB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err = db.Close(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Postgres connection closed")
}
