package dump

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Credits struct {
	Directors []string `xml:"director"` // epg: directed
	Actors    []string `xml:"actor"`    // epg: actors
	Producers []string `xml:"producer"` // epg: += directed
	Composers []string `xml:"composer"` // epg: += directed
	Writers   []string `xml:"writer"`   // epg: += directed
}

type Rating struct {
	Value int `xml:"value"`
}

type ProgInfo struct {
	Start   string `xml:"start,attr"`   // epg: utc_date, epg_date = utc_date + 3h,   utc_start, date_start = utc_date + 3h
	Stop    string `xml:"stop,attr"`    // epg:                                       utc_stop, date_stop = utc_stop + 3h
	Channel string `xml:"channel,attr"` // epg_mapping: epg_code
	ChEpgId int

	Title      string   `xml:"title"` // epg: title
	SubTitle   string   `xml:"sub-title"`
	Desc       string   `xml:"desc"`     // epg: description = subtitle + description
	Categories []string `xml:"category"` // epg: genres join(,)
	Credits    Credits  `xml:"credits"`
	Date       string   `xml:"date"`    // .start.date
	Countries  []string `xml:"country"` // epg: country join(,)
	Year       string   `xml:"year"`    // epg: create_year
	Rating     Rating   `xml:"rating"`  // epg: minage
}

type Tv struct {
	Programme []ProgInfo `xml:"programme"`
}

func (tv *Tv) ByChanCode(epgCode string) []ProgInfo {
	res := make([]ProgInfo, 0)
	for _, pInfo := range tv.Programme {
		if pInfo.Channel == epgCode {
			res = append(res, pInfo)
		}
	}
	return res
}

func (pInf *ProgInfo) GetTimes() (start time.Time, stop time.Time, err error) {
	start, err = time.Parse("20060102150405 -0700", pInf.Start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	stop, err = time.Parse("20060102150405 -0700", pInf.Stop)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return
}

func Parse(dumpPath string) (*Tv, error) {
	b, err := ioutil.ReadFile(filepath.Clean(dumpPath))
	if err != nil {
		return nil, err
	}

	var xmlTv Tv

	if err = xml.Unmarshal(b, &xmlTv); err != nil {
		return nil, err
	}

	return &xmlTv, nil
}
