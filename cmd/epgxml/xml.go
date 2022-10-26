package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type ProgInfo struct {
	Start   string `xml:"start,attr"`   // epg: utc_date, epg_date = utc_date + 3h,   utc_start, date_start = utc_date + 3h
	Stop    string `xml:"stop,attr"`    // epg:                                       utc_stop, date_stop = utc_stop + 3h
	Channel string `xml:"channel,attr"` // epg_mapping: epg_code
	ChEpgId int

	Title      string   `xml:"title"` // epg: title
	SubTitle   string   `xml:"sub-title"`
	Desc       string   `xml:"desc"`     // epg: description = subtitle + description
	Categories []string `xml:"category"` // epg: genres join(,)
	Credits    struct {
		Directors []string `xml:"director"` // epg: directed
		Actors    []string `xml:"actor"`    // epg: actors
		Producers []string `xml:"producer"` // epg: += directed
		Composers []string `xml:"composer"` // epg: += directed
		Writers   []string `xml:"writer"`   // epg: += directed
	} `xml:"credits"`
	Date      string   `xml:"date"`    // .start.date
	Countries []string `xml:"country"` // epg: country join(,)
	Year      string   `xml:"year"`    // epg: create_year
	Rating    struct {
		Value int `xml:"value"`
	} `xml:"rating"` // epg: minage
}

type Tv struct {
	Programme []ProgInfo `xml:"programme"`
}

func (tv *Tv) getByChannel(epgCode string) []ProgInfo {
	res := make([]ProgInfo, 0)
	for _, pInfo := range tv.Programme {
		if pInfo.Channel == epgCode {
			res = append(res, pInfo)
		}
	}
	return res
}

func readXml(xmlFilePath string) (*Tv, error) {
	fnFile, err := os.Open(xmlFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := fnFile.Close(); err != nil {
			LOG.Panic(err)
		}
	}()

	xmlRaw, _ := ioutil.ReadAll(fnFile)
	var xmlTv Tv

	if err = xml.Unmarshal(xmlRaw, &xmlTv); err != nil {
		return nil, err
	}

	return &xmlTv, nil
}
