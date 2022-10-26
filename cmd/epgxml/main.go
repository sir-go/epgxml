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
		LOG.Panic(err)
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
	LOG.Info("-- start --")

	LOG.Debug("process command line arguments")
	fCfgPath := flag.String("c", "conf.toml", "path to conf file")
	flag.Parse()

	LOG.Debug("parse config file")
	cfg, err := LoadConfig(*fCfgPath)
	check(err)

	LOG.Debugf("connect to db %s", cfg.Db.Dbpath)
	dbConn, err := dbConnect(cfg.Db.User, cfg.Db.Password, cfg.Db.Dbpath)
	check(err)
	defer func() {
		if err := dbConn.Close(); err != nil {
			LOG.Panic(err)
		}
	}()

	LOG.Debug("get channels from epg_mapping")
	epgMapping, err := getEpgMapping(dbConn)
	check(err)
	if len(epgMapping) < 1 {
		LOG.Debug("table epg_mapping is empty")
		return
	}
	LOG.Debugf("got %d records", len(epgMapping))

	LOG.Debugf("get programme from %s", cfg.Xml.FileName)
	xmlTv, err := readXml(cfg.Xml.FileName)
	check(err)
	if len(xmlTv.Programme) < 1 {
		LOG.Debugf("%s has no programme records", cfg.Xml.FileName)
		return
	}
	LOG.Debugf("got %d records", len(xmlTv.Programme))

	LOG.Debug("clear epg table")
	check(clearEpg(dbConn))
	pb := progressbar.New(len(epgMapping))

	dbTx, err := dbConn.Begin()
	check(err)
	defer func() {
		if err := dbTx.Rollback(); err != nil {
			LOG.Panic(err)
		}
	}()

	LOG.Debugf("xml -> db convert started")
	for _, chRec := range epgMapping {
		for _, progRec := range xmlTv.getByChannel(chRec.EpgCode) {
			timeStart, err := time.Parse("20060102150405 -0700", progRec.Start)
			if err != nil {
				LOG.Warningf("can't parse start time from %s", timeStart)
				continue
			}

			timeStop, err := time.Parse("20060102150405 -0700", progRec.Stop)
			if err != nil {
				LOG.Warningf("can't parse stop time from %s", timeStop)
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
			// LOG.Debug(eRec)
			check(addEpgRecord(dbTx, eRec))
		}
		if err := pb.Add(1); err != nil {
			LOG.Panic(err)
		}
	}
	if _, err = os.Stdout.WriteString("\n"); err != nil {
		LOG.Panic(err)
	}
	LOG.Debug("commit db changes")
	if err = dbTx.Commit(); err != nil {
		LOG.Panic(err)
	}

	LOG.Debug("update epg_updated in dvb_network and dvb_streams")
	check(updDates(dbConn))
	LOG.Info("-- done --")
}
