package utils

import (
	"log"
	"math"
	"strings"

	"epgxml/internal/db"
	"epgxml/internal/dump"
)

func dumpRecordToDbRecord(chanId int, pInfo dump.ProgInfo) (dbRec *db.Record) {
	timeStart, timeStop, err := pInfo.GetTimes()
	if err != nil {
		log.Println("start and stop time parsing", err)
		return nil
	}

	dbRec = &db.Record{
		ChId:        chanId,
		EpgDate:     timeStart.Local(),
		UtcDate:     timeStart,
		DateStart:   timeStart.Local(),
		DateStop:    timeStop.Local(),
		UtcStart:    timeStart,
		UtcStop:     timeStop,
		Duration:    int(math.Round(timeStop.Sub(timeStart).Minutes())),
		Title:       CutStr(pInfo.Title, 255),
		Description: CutStr(strings.Join([]string{pInfo.SubTitle, pInfo.Desc}, " "), 4096),
		Genres:      CutStr(strings.Join(pInfo.Categories, ", "), 255),
		MinAge:      pInfo.Rating.Value,
		CreateYear:  CutStr(pInfo.Year, 255),
		Actors:      CutStr(strings.Join(pInfo.Credits.Actors, ", "), 255),
		Directed:    CutStr(strings.Join(pInfo.Credits.Directors, ", "), 255),
		Country:     CutStr(strings.Join(pInfo.Countries, ", "), 255)}
	return
}

func GenDbRecords(chanId int, chanDump []dump.ProgInfo) (dbRecords []db.Record) {
	dbRecords = make([]db.Record, 0)
	for _, progRec := range chanDump {
		dbRecord := dumpRecordToDbRecord(chanId, progRec)
		if dbRecord == nil {
			continue
		}
		dbRecords = append(dbRecords, *dbRecord)
	}
	return
}
