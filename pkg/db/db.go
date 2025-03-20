package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/ssh"
)

func WithTransaction(db *bun.DB, fn func(tx bun.Tx) error, opts *sql.TxOptions, c domain.Context) error {
	tx, err := db.BeginTx(c.Request().Context(), opts)

	if err != nil {
		return err
	}

	err = fn(tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// JoinSuffix returns the `ON` part of a `JOIN` statement for two table that shares the same key name
//
// example:
// JoinSuffix("companies", "projects", "company_id") => "companies ON companies.company_id = projects.company_id"
func JoinSuffix(table1, table2, column string) string {
	return table1 + " ON " + table1 + "." + column + " = " + table2 + "." + column
}

func QueryExecuteError(err error, query string, args ...interface{}) error {
	return fmt.Errorf("failed to execute query: %w\n%s\n%s", err, query, args)
}

func UnmarshalError(err error, query string, args []interface{}) error {
	return fmt.Errorf("failed to unmarshal struct: %w\n%s\n%s", err, query, args)
}

// PostgresDSN returns a connection URL for postgres
func PostgresDSN(host, port, user, pass, database string) string {
	// postgres wants connection url to be percentage quoted
	// https://www.postgresql.org/docs/11/libpq-connect.html#id-1.7.3.8.3.6
	pass = url.QueryEscape(pass)

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, database)
}

// Postgres creates a database connection with pgx and sqlx
func Postgres(host, port, user, pass, database string, clients ...*ssh.Client) (*bun.DB, error) {
	dsn := PostgresDSN(host, port, user, pass, database)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	return bun.NewDB(sqldb, pgdialect.New()), nil
}

func WithRetry(period time.Duration, limit int, fn func() (*bun.DB, error)) (*bun.DB, error) {
	var db *bun.DB
	var err error

	for i := 0; i < limit; i++ {
		if db, err = fn(); err == nil {
			return db, nil
		}
		<-time.After(period)
	}

	return nil, fmt.Errorf("retry failed: %w", err)
}

func SSH(sshUser, sshHost, sshPort string, sshKey []byte) (*ssh.Client, error) {
	key, err := ssh.ParsePrivateKey(sshKey)

	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	auth := ssh.PublicKeys(key)

	cfg := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshHost, sshPort), cfg)

	if err != nil {
		return nil, fmt.Errorf("unable to open ssh connection: %w", err)
	}

	return client, nil
}

// GetDBColumns returns all the field defined by `db` tags of a struct
func GetDBColumns(v interface{}) []string {
	t := reflect.TypeOf(v)
	var columns []string

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if column, ok := f.Tag.Lookup("db"); ok {
			columns = append(columns, column)
		}
	}

	return columns
}
