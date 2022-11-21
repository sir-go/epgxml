package tests

import (
	"reflect"
	"testing"

	"epgxml/internal/db"
	"epgxml/internal/dump"
	"epgxml/internal/utils"
)

func TestCutStr(t *testing.T) {
	type args struct {
		s      string
		maxLen int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{"", 0}, ""},
		{"less", args{"abcd", 5}, "abcd"},
		{"more", args{"abcdefg", 5}, "abcde"},
		{"eq", args{"abcd", 4}, "abcd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.CutStr(tt.args.s, tt.args.maxLen); got != tt.want {
				t.Errorf("CutStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenDbRecords(t *testing.T) {
	type args struct {
		chanId   int
		chanDump []dump.ProgInfo
	}
	tests := []struct {
		name          string
		args          args
		wantDbRecords []db.Record
	}{
		{"empty", args{0,
			[]dump.ProgInfo{}},
			[]db.Record{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDbRecords := utils.GenDbRecords(tt.args.chanId, tt.args.chanDump)
			if !reflect.DeepEqual(gotDbRecords, tt.wantDbRecords) {
				t.Errorf("GenDbRecords() = %v, want %v", gotDbRecords, tt.wantDbRecords)
			}
		})
	}
}
