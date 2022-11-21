package tests

import (
	"encoding/gob"
	"os"
	"reflect"
	"testing"
	"time"

	"epgxml/internal/dump"
)

func loadDump() *dump.Tv {
	fd, err := os.Open("testdata/TV_Pack.gob")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}()
	d := &dump.Tv{}
	if err = gob.NewDecoder(fd).Decode(d); err != nil {
		panic(err)
	}
	return d
}

func TestParse(t *testing.T) {
	d := loadDump()

	type args struct {
		dumpPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *dump.Tv
		wantErr bool
	}{
		{"empty", args{""}, nil, true},
		{"ok", args{"testdata/TV_Pack.xml"}, d, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dump.Parse(tt.args.dumpPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgInfo_GetTimes(t *testing.T) {
	d := loadDump()

	tests := []struct {
		name      string
		pInf      dump.ProgInfo
		wantStart time.Time
		wantStop  time.Time
		wantErr   bool
	}{
		{"empty", dump.ProgInfo{}, time.Time{}, time.Time{}, true},
		{"ok",
			d.Programme[0],
			time.Unix(1534107600, 0).UTC(),
			time.Unix(1534110600, 0).UTC(),
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotStop, err := tt.pInf.GetTimes()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTimes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !gotStart.Equal(tt.wantStart) {
				t.Errorf("GetTimes() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if !gotStop.Equal(tt.wantStop) {
				t.Errorf("GetTimes() gotStop = %v, want %v", gotStop, tt.wantStop)
			}
		})
	}
}

func TestTv_ByChanCode(t *testing.T) {
	d := loadDump()

	type args struct {
		epgCode string
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
	}{
		{"empty", args{""}, 0},
		{"ok", args{"000000048"}, 775},
		{"ok", args{"000000206"}, 217},
		{"цкщтп", args{"00000999"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := d.ByChanCode(tt.args.epgCode); len(got) != tt.wantLen {
				t.Errorf("ByChanCode() length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}
