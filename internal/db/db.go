package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/nakagami/firebirdsql"
)

type EpgDB struct {
	dsnString string
}

type Mapping struct {
	ChId    int    `db:"CH_ID"`
	EpgCode string `db:"EPG_CODE"`
}

//goland:noinspection SpellCheckingInspection
type Record struct {
	ChId        int       `db:"CH_ID"`
	EpgDate     time.Time `db:"EPG_DATE"`
	UtcDate     time.Time `db:"UTC_DATE"`
	DateStart   time.Time `db:"DATE_START"`
	DateStop    time.Time `db:"DATE_STOP"`
	UtcStart    time.Time `db:"UTC_START"`
	UtcStop     time.Time `db:"UTC_STOP"`
	Duration    int       `db:"DURATION"`
	Title       string    `db:"TITLE"`
	Description string    `db:"DESCRIPTION"`
	Genres      string    `db:"GENRES"`
	MinAge      int       `db:"MINAGE"`
	CreateYear  string    `db:"CREATE_YEAR"`
	Actors      string    `db:"ACTORS"`
	Directed    string    `db:"DIRECTED"`
	Country     string    `db:"COUNTRY"`
}

// FormatDSN gets db connection arguments and returns formatted DSN string
func FormatDSN(host string, port int, username, password, dbPath string) string {
	return fmt.Sprintf("%s:%s@%s:%d/%s?lc_ctype=utf8",
		username, password, host, port, dbPath)
}

// New creates new EpgDB structure
//goland:noinspection GoUnusedExportedFunction
func New(dsnString string) *EpgDB {
	return &EpgDB{dsnString: dsnString}
}

// WithTx is a wrapper around db transaction
func (d *EpgDB) WithTx(do func(tx *sql.Tx) error) error {
	conn, err := sql.Open("firebirdsql", d.dsnString)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("db connection closing", err)
		}
	}()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	if err = do(tx); err != nil {
		defer func() {
			if err := tx.Rollback(); err != nil {
				log.Println("db transaction rollback", err)
			}
		}()
		return err
	}
	return nil
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection
func (d *EpgDB) GetMapping() (res []Mapping, err error) {
	res = make([]Mapping, 0)
	err = d.WithTx(func(tx *sql.Tx) error {
		rows, err := tx.Query("select CH_ID, EPG_CODE from EPG_MAPPING")
		if err != nil {
			return err
		}

		for rows.Next() {
			var record Mapping
			if err := rows.Scan(&record.ChId, &record.EpgCode); err != nil {
				return err
			}
			res = append(res, record)
		}
		if err := rows.Err(); err != nil {
			return err
		}

		return nil
	})
	return
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection,SqlWithoutWhere
func (d *EpgDB) Clear() error {
	return d.WithTx(func(tx *sql.Tx) error {
		_, err := tx.Exec("delete from epg")
		if err != nil {
			return err
		}
		return tx.Commit()
	})
}

func (d *EpgDB) AddMany(records []Record) error {
	//goland:noinspection SpellCheckingInspection,SqlResolve,SqlNoDataSourceInspection
	const qry = `insert into epg (
	CH_ID, EPG_DATE, UTC_DATE, DATE_START, DATE_STOP, 
	UTC_START, UTC_STOP, DURATION, TITLE, DESCRIPTION, GENRES, MINAGE, 
	CREATE_YEAR, ACTORS, DIRECTED, COUNTRY) values (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return d.WithTx(func(tx *sql.Tx) error {
		for _, r := range records {
			_, err := tx.Exec(qry,
				r.ChId, r.EpgDate, r.UtcDate, r.DateStart, r.DateStop,
				r.UtcStart, r.UtcStop, r.Duration, r.Title, r.Description,
				r.Genres, r.MinAge, r.CreateYear, r.Actors, r.Directed, r.Country)
			if err != nil {
				return err
			}
		}
		return tx.Commit()
	})
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection,SqlWithoutWhere
func (d *EpgDB) Touch() error {
	return d.WithTx(func(tx *sql.Tx) (err error) {
		_, err = tx.Exec(
			"update dvb_network set epg_updated = current_timestamp")
		if err != nil {
			return err
		}
		_, err = tx.Exec(
			"update dvb_streams set epg_updated = current_timestamp")
		if err != nil {
			return err
		}
		return tx.Commit()
	})
}
