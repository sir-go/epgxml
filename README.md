# EPG converter from XML dump to local Firebird Database
[![Go](https://github.com/sir-go/epgxml/actions/workflows/go.yml/badge.svg)](https://github.com/sir-go/epgxml/actions/workflows/go.yml)

## What it does
- script `get_xml.sh` downloads an EPG XML dump from the EPG service site,
  it uses a `ncftpget` tool from the `ncftp` package
- the app parses an XML dump to the local Firebird database
- updates the dates in the tables `dvb_network` and `dvb_streams`

## Configuration
All the settings are set in the `config.yml` file (path to config file can be set
by `-c` option)

```yaml
username:   sysdba                    # database username
password:   masterkey                 # database password
db_path:    /firebird/data/a4on.fdb   # database file location
dump_path:  /TV_Pack.xml              # downloaded XML dump locatin
host:       firebird-db-test          # firebird server host
port:       3050                      # firebird server port

```
## Tests
You should run testing firebird server instance before running tests
```bash
docker compose up -d

go test -v ./tests
gosec ./...

docker compose down
```

## Docker
```bash
docker build -t epgxml .

docker run \
  --net epgxml_net \
  --name epgxml \
  --rm -it \
  -v ${PWD}/tests/testdata/a4on.fdb:/firebird/data/a4on.fdb \
  -v ${PWD}/tests/testdata/TV_Pack.xml:/TV_Pack.xml \
  -v ${PWD}/config.yml:/config.yml \
  epgxml:latest
```

## Standalone
```bash
go mod download
go build -o epgxml ./cmd/epgxml

./epgxml -c config.yml
```
