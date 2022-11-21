package tests

import (
	"database/sql"
	"testing"

	"epgxml/internal/db"
)

// alert: run local test fb3 server before

const FbPort = 43050

func TestFormatDSN(t *testing.T) {
	type args struct {
		host     string
		port     int
		username string
		password string
		dbPath   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty",
			args{"", 0, "", "", ""},
			":@:0/?lc_ctype=utf8"},
		{"ok",
			args{"localhost", 3050, "uname",
				"passwd", "/path/to/db.fdb"},
			"uname:passwd@localhost:3050//path/to/db.fdb?lc_ctype=utf8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.FormatDSN(
				tt.args.host, tt.args.port, tt.args.username,
				tt.args.password, tt.args.dbPath); got != tt.want {
				t.Errorf("FormatDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setUpConnection() *db.EpgDB {
	return db.New(
		db.FormatDSN(
			"localhost", FbPort, "sysdba",
			"masterkey", "/firebird/data/a4on.fdb"))
}

func TestEpgDB_Clear(t *testing.T) {
	d := setUpConnection()
	if err := d.AddMany([]db.Record{
		{ChId: 990, Title: "some-title-0"},
		{ChId: 991, Title: "some-title-1"},
		{ChId: 992, Title: "some-title-2"},
		{ChId: 993, Title: "some-title-3"},
		{ChId: 994, Title: "some-title-4"},
	}); err != nil {
		panic(err)
	}
	defer func() {
		if err := d.Clear(); err != nil {
			panic(err)
		}
	}()

	t.Run("test clear", func(t *testing.T) {
		if err := d.Clear(); (err != nil) != false {
			t.Errorf("Clear() error = %v, wantErr %v", err, false)
		}
		c, err := getRecordsCount(d)
		if err != nil {
			t.Errorf("getRecordsCount error %v", err)
		}
		if c != 0 {
			t.Errorf("Clear() must delete all receords, by still %d", c)
		}
	})
}

func getRecordsCount(d *db.EpgDB) (count int, err error) {
	err = d.WithTx(func(tx *sql.Tx) error {
		//goland:noinspection SqlResolve
		rows, err := tx.Query("select count(*) from EPG")
		if err != nil {
			return err
		}

		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				return err
			}
		}
		if err = rows.Err(); err != nil {
			return err
		}
		return nil
	})
	return
}

func TestEpgDB_AddMany(t *testing.T) {
	d := setUpConnection()
	defer func() {
		if err := d.Clear(); err != nil {
			panic(err)
		}
	}()
	tests := []struct {
		name             string
		records          []db.Record
		wantRecordsCount int
		wantErr          bool
	}{
		{"empty", []db.Record{}, 0, false},
		{"some", []db.Record{
			{ChId: 990, Title: "some-title-0"},
			{ChId: 991, Title: "some-title-1"},
			{ChId: 992, Title: "some-title-2"},
			{ChId: 993, Title: "some-title-3"},
			{ChId: 994, Title: "some-title-4"},
		}, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.AddMany(tt.records); (err != nil) != tt.wantErr {
				t.Errorf("AddMany() error = %v, wantErr %v", err, tt.wantErr)
			}
			c, err := getRecordsCount(d)
			if err != nil {
				t.Errorf("getRecordsCount error %v", err)
			}
			if c != tt.wantRecordsCount {
				t.Errorf("getRecordsCount() = %d, want %d", c, tt.wantRecordsCount)
			}
		})
	}
}

func TestEpgDB_GetMapping(t *testing.T) {
	d := setUpConnection()
	t.Run("get-some", func(t *testing.T) {
		gotRes, err := d.GetMapping()
		if err != nil {
			t.Errorf("GetMapping() error = %v, wantErr %v", err, false)
			return
		}
		if len(gotRes) < 1 {
			t.Error("GetMapping() length < 1")
		}
	})
}

func TestEpgDB_Touch(t *testing.T) {
	d := setUpConnection()
	t.Run("", func(t *testing.T) {
		err := d.Touch()
		if err != nil {
			t.Errorf("Touch() error = %v, wantErr %v", err, false)
		}
	})
}
