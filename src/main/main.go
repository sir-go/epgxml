package main

import (
	"flag"
	"math"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func init() {
	initLogging("epg-xml")
}

func cutStr(s string, maxLen int) string {
	r := []rune(s)
	if len(r) > maxLen {
		return string(r[:maxLen])
	} else {
		return s

	}
}

func main() {
	log.Info("-- start --")

	log.Debug("process command line arguments")
	fCfgPath := flag.String("c", "conf.toml", "path to conf file")
	flag.Parse()

	log.Debug("parse config file")
	cfg, err := LoadConfig(*fCfgPath)
	check(err)

	log.Debugf("connect to db %s", cfg.Db.Dbpath)
	dbConn, err := dbConnect(cfg.Db.User, cfg.Db.Password, cfg.Db.Dbpath)
	check(err)
	defer dbConn.Close()

	log.Debug("get channels from epg_mapping")
	epgMapping, err := getEpgMapping(dbConn)
	check(err)
	if len(epgMapping) < 1 {
		log.Debug("table epg_mapping is empty")
		return
	}
	log.Debugf("got %d records", len(epgMapping))

	log.Debugf("get programme from %s", cfg.Xml.FileName)
	xmlTv, err := readXml(cfg.Xml.FileName)
	check(err)
	if len(xmlTv.Programme) < 1 {
		log.Debugf("%s has no programme records", cfg.Xml.FileName)
		return
	}
	log.Debugf("got %d records", len(xmlTv.Programme))

	log.Debug("clear epg table")
	check(clearEpg(dbConn))
	pb := progressbar.New(len(epgMapping))

	dbTx, err := dbConn.Begin()
	check(err)
	defer dbTx.Rollback()

	log.Debugf("xml -> db convert started")
	for _, chRec := range epgMapping {
		for _, progRec := range xmlTv.getByChannel(chRec.EpgCode) {
			timeStart, err := time.Parse("20060102150405 -0700", progRec.Start)
			if err != nil {
				log.Warningf("can't parse start time from %s", timeStart)
				continue
			}

			timeStop, err := time.Parse("20060102150405 -0700", progRec.Stop)
			if err != nil {
				log.Warningf("can't parse stop time from %s", timeStop)
				continue
			}

			eRec := &Epg{
				ChId:        chRec.ChId,
				EpgDate:     timeStart.Local(),
				UtcDate:     timeStart,
				DateStart:   timeStart.Local(),
				DateStop:    timeStop.Local(),
				UtcStart:    timeStart,
				UtcStop:     timeStop,
				Duration:    int(math.Round(timeStop.Sub(timeStart).Minutes())),
				Title:       cutStr(progRec.Title, 255),
				Description: cutStr(strings.Join([]string{progRec.SubTitle, progRec.Desc}, " "), 4096),
				Genres:      cutStr(strings.Join(progRec.Categories, ", "), 255),
				MinAge:      progRec.Rating.Value,
				CreateYear:  cutStr(progRec.Year, 255),
				Actors:      cutStr(strings.Join(progRec.Credits.Actors, ", "), 255),
				Directed:    cutStr(strings.Join(progRec.Credits.Directors, ", "), 255),
				Country:     cutStr(strings.Join(progRec.Countries, ", "), 255)}
			// log.Debug(eRec)
			check(addEpgRecord(dbTx, eRec))
		}
		pb.Add(1)
	}
	os.Stdout.WriteString("\n")
	log.Debug("commit db changes")
	dbTx.Commit()
	pb.Finish()

	log.Debug("update epg_updated in dvb_network and dvb_streams")
	check(updDates(dbConn))
	log.Info("-- done --")
}
