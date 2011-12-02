package replay

import (
	"os"
	"io"
	"io/ioutil"
	"json"
	"log"
	"bytes"
	"compress/gzip"
	"regexp"
)

var htmlReplayRe *regexp.Regexp = regexp.MustCompile("loadReplayData\\([:space]*'[:space:]*({.+})[:space:]*'[:space:]*\\)[:space:]*;")

// Take a job and load the required data
func Load(file string) (*Match, os.Error) {
	m := Match{}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// decode gzip file if needed
	if buf[0] == '\x1f' && buf[1] == '\x8b' {
		unzip, err := gzip.NewReader(bytes.NewBuffer(buf[:]))
		if err != nil {
			log.Panic("Unzip: %v", err)
		}
		bout := bytes.NewBuffer(make([]byte, 0, 8*len(buf)))
		_, err = io.Copy(bout, unzip)
		if err != nil {
			log.Panic("Unzip: %v", err)
		}
		buf = bout.Bytes()
	}

	// extract replay if html
	match := htmlReplayRe.FindSubmatchIndex(buf)
	if match != nil {
		//log.Printf("Matches %v", match)
		//log.Printf("Matches %s .. %s", buf[match[0]:match[0]+15], buf[match[1]-15:match[1]])
		//log.Printf("Matches %s .. %s", buf[match[2]:match[2]+15], buf[match[3]-15:match[3]])
		buf = buf[match[2]:match[3]]
	}

	err = json.Unmarshal(buf, &m)
	// TODO hackery - rows and cols belong in the map not the gameinfo
	m.Rows = m.Map.Rows
	m.Cols = m.Map.Cols
	m.Replay.GameInfo.PlayerSeed = m.Replay.PlayerSeed

	if err != nil {
		return nil, err
	}

	return &m, nil
}
