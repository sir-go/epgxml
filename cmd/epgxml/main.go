package main

import (
	"flag"
	"log"

	"github.com/schollz/progressbar"

	"epgxml/internal/config"
	"epgxml/internal/db"
	"epgxml/internal/dump"
	"epgxml/internal/utils"
)

func main() {
	log.Println("-- start --")
	cfgPath := flag.String("c", "config.yml",
		"path to config file")
	flag.Parse()

	log.Println("fill the config")
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		panic(err)
	}

	epgDb := db.New(
		db.FormatDSN(
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbPath))

	log.Println("get channels from epg_mapping")
	epgMapping, err := epgDb.GetMapping()
	if err != nil {
		panic(err)
	}

	if len(epgMapping) < 1 {
		log.Println("table epg_mapping is empty")
		return
	}
	log.Println("got records:", len(epgMapping))

	log.Println("get programme from", cfg.DumpPath)
	xmlTv, err := dump.Parse(cfg.DumpPath)
	if err != nil {
		panic(err)
	}

	if len(xmlTv.Programme) < 1 {
		log.Println(cfg.DumpPath, "has no programme records")
		return
	}
	log.Println("got records:", len(xmlTv.Programme))

	log.Println("clear epg table")
	if err != nil {
		panic(err)
	}

	pb := progressbar.New(len(epgMapping))

	log.Println("xml -> db convert started")
	dbRecords := make([]db.Record, 0)
	for _, chRec := range epgMapping {
		dbRecords = append(dbRecords,
			utils.GenDbRecords(chRec.ChId, xmlTv.ByChanCode(chRec.EpgCode))...)
		if err := pb.Add(1); err != nil {
			log.Println(err)
		}
	}
	log.Println("done")

	log.Println("insert data to db...")
	if err = epgDb.AddMany(dbRecords); err != nil {
		panic(err)
	}

	log.Println("update epg_updated in dvb_network and dvb_streams")
	if err = epgDb.Touch(); err != nil {
		log.Println(err)
	}

	log.Println("-- done --")
}
