## EPG converter from XML dump to local Firebird Database

### What it does

- script `get_xml.sh` downloads an EPG XML dump from the EPG service site,
  it uses a `ncftpget` tool from the `ncftp` package
- the app parses an XML dump to the local Firebird database
- updates the dates in the tables `dvb_network` and `dvb_streams`
___
### Configuration

All the settings are set in the `conf.toml` file (path to config file can be set
by `-c` option)

```toml
[db]                    # local database connection settings                      
  user = ""             # db user
  password = ""         # db password
  dbpath = "/var/lib/firebird/2.5/data/A4ON_FREE.FDB" # database file location

[xml]
    filename = "TV_Pack.xml"    # gotten XML dump location
```
___
### Build and run
```bash
go mod download
build -o epgxml ./cmd/epgxml
```
