package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/nakagami/firebirdsql"
)

type EpgMapping struct {
	ChId    int    `db:"CH_ID"`
	EpgCode string `db:"EPG_CODE"`
}

type Epg struct {
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

func dbConnect(user string, password string, dbpath string) (*sql.DB, error) {
	return sql.Open("firebirdsql",
		fmt.Sprintf("%s:%s@localhost%s", user, password, dbpath))
}

func getEpgMapping(dbConn *sql.DB) ([]EpgMapping, error) {
	rows, err := dbConn.Query("select CH_ID, EPG_CODE from EPG_MAPPING")
	if err != nil {
		return nil, err
	}

	res := make([]EpgMapping, 0)

	for rows.Next() {
		var record EpgMapping
		if err := rows.Scan(&record.ChId, &record.EpgCode); err != nil {
			return nil, err
		}
		res = append(res, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func clearEpg(dbConn *sql.DB) error {
	dbTx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := dbTx.Rollback(); err != nil {
			LOG.Panic(err)
		}
	}()

	_, err = dbTx.Exec("delete from epg")
	if err != nil {
		return err
	}

	return dbTx.Commit()
}

func addEpgRecord(dbTx *sql.Tx, r *Epg) error {
	//noinspection ALL
	qry := `insert into epg (CH_ID, EPG_DATE, UTC_DATE, DATE_START, DATE_STOP, UTC_START, UTC_STOP, DURATION, TITLE, 
DESCRIPTION, GENRES, MINAGE, CREATE_YEAR, ACTORS, DIRECTED, COUNTRY) values (
?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := dbTx.Exec(qry, r.ChId, r.EpgDate, r.UtcDate, r.DateStart, r.DateStop, r.UtcStart, r.UtcStop,
		r.Duration, r.Title, r.Description, r.Genres, r.MinAge, r.CreateYear, r.Actors, r.Directed, r.Country)
	return err
}

func updDates(dbConn *sql.DB) error {
	if _, err := dbConn.Exec("update dvb_network set epg_updated = current_timestamp"); err != nil {
		return err
	}

	if _, err := dbConn.Exec("update dvb_streams set epg_updated = current_timestamp"); err != nil {
		return err
	}

	return nil
}
