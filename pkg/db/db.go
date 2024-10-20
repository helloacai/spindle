package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/helloacai/spindle/pkg/log"
)

var table string

func init() {
	table = os.Getenv("CURSOR_TABLE_NAME")
}

func SaveCursor(cursor string) error {
	log.Debug().Str("cursor", cursor).Msg("saving cursor")
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Err(err).Msg("Unable to connect to database")
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "update "+table+" set cursor = $1", cursor)
	if err != nil {
		log.Err(err).Msg("Exec failed")
		return err
	}

	log.Debug().Msg("cursor saved")
	return nil
}

func GetCursor() (string, error) {
	log.Debug().Msg("fetching cursor")
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Err(err).Msg("Unable to connect to database")
		return "", err
	}
	defer conn.Close(context.Background())

	var cursor string
	err = conn.QueryRow(context.Background(), "select cursor from "+table+" LIMIT 1").Scan(&cursor)
	if err != nil {
		log.Err(err).Msg("QueryRow failed")
		return "", err
	}

	log.Debug().Str("cursor", cursor).Msg("cursor fetched")
	return cursor, nil
}
