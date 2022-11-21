package tests

import (
	"reflect"
	"testing"

	"epgxml/internal/config"
)

func TestLoad(t *testing.T) {
	type args struct {
		confPath string
	}
	tests := []struct {
		name    string
		args    args
		wantCfg *config.Config
		wantErr bool
	}{
		{"empty", args{""}, nil, true},
		{"404", args{"/fake/path"}, nil, true},
		{"bad", args{"testdata/config-bad.yml"},
			&config.Config{
				Password: "password-value",
				DumpPath: "path-to-dump-xml",
			}, true},
		{"good", args{"testdata/config.yml"},
			&config.Config{
				Username: "username-value",
				Password: "password-value",
				DbPath:   "path-to-database",
				DumpPath: "path-to-dump-xml",
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := config.Load(tt.args.confPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("Load() gotCfg = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}
}
