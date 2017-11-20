package mp3lib

import (
	"os"
	"strconv"

	"io/ioutil"
	"sort"

	id3 "id3-go"
	"path/filepath"
	"strings"
)

func (i *MP3Library) SetRating(fname string, rating int) int {
	var err error

	i.mutex.Lock()
	defer i.mutex.Unlock()

	fullPath := fname
	if !strings.HasPrefix(fname, "/") {
		fullPath, err = filepath.Abs(i.BaseDir + "/" + fname)
		if err != nil {
			log.Warningf("SetRating filepath.Abs: %v", err)
		}
	}

	fd, err := id3.Open(fullPath)
	if err != nil {
		log.Warningf("MP3Library.SetRating id3.Open: %v", err)
		return RATING_UNKNOWN
	}
	defer fd.Close()

	newRating := strconv.Itoa(rating)
	fd.SetYear(newRating)

	log.Infof("MP3Library.SetRating: Rating for %s set to %d", fname, RATING_DEFAULT)

	return rating
}

func (i *MP3Library) GetRating(fname string) int {
	var err error

	fullPath := fname
	if !strings.HasPrefix(fname, "/") {
		fullPath, err = filepath.Abs(i.BaseDir + "/" + fname)
		if err != nil {
			log.Warningf("SetRating filepath.Abs: %v", err)
		}
	}

	if _, err := os.Stat(fullPath); err != nil {
		log.Warningf("MP3Library.GetRating os.Stat: %v", err)
		return RATING_UNKNOWN
	}

	i.mutex.Lock()
	defer i.mutex.Unlock()

	fd, err := id3.Open(fullPath)
	if err != nil {
		log.Warningf("MP3Library.GetRating id3.Open: %v", err)
		return RATING_UNKNOWN
	}
	defer fd.Close()

	curRating_s := fd.Year()

	if curRating_s == "" {
		return RATING_UNKNOWN
	}

	curRating, err := strconv.Atoi(curRating_s)
	if err != nil {
		log.Warningf("MP3Library.GetRating strconv.Atoi: %v", err)
		return RATING_UNKNOWN
	}

	return curRating
}

func (i *MP3Library) DecreaseRating(name string) int {
	curRating := i.GetRating(name)

	switch curRating {
	case RATING_UNKNOWN:
		return RATING_UNKNOWN
	case RATING_ZERO:
		return RATING_ZERO
	default:
		{
			curRating -= 1

			rating := i.SetRating(name, curRating)

			return rating
		}
	}

	return RATING_UNKNOWN
}

func (i *MP3Library) IncreaseRating(name string) int {
	curRating := i.GetRating(name)

	switch curRating {
	case RATING_UNKNOWN:
		return RATING_UNKNOWN
	case RATING_MAX:
		return RATING_MAX
	default:
		{
			curRating += 1

			rating := i.SetRating(name, curRating)

			return rating
		}

	}

	return RATING_UNKNOWN
}

func (i *MP3Library) RemoveFile(fname string) bool {
	var err error

	fullPath := fname
	if !strings.HasPrefix(fname, "/") {
		fullPath, err = filepath.Abs(i.BaseDir + "/" + fname)
		if err != nil {
			log.Warningf("SetRating filepath.Abs: %v", err)
		}
	}

	_, err = os.Stat(fullPath)
	if err != nil {
		log.Warningf("MP3Library.RemoveFile os.Stat: %v", err)
		return false
	}

	err = os.Remove(fullPath)
	if err != nil {
		log.Warningf("MP3Library.RemoveFile os.Remove: %v", err)
		return false
	}

	return true
}

func (i *MP3Library) GetAllFiles() []string {

	files, err := ioutil.ReadDir(i.BaseDir)
	if err != nil {
		log.Warningf("MP3Library.GetAllFiles ioutil.ReadDir: %v", err)
		return nil
	}

	tmpList := make([]string, len(files))
	totItems := 0
	for _, fs := range files {
		if fs.Name() == "" {
			continue
		}
		tmpList = append(tmpList, fs.Name())
		totItems += 1
	}

	response := make([]string, totItems)
	response = tmpList

	sort.Strings(response)

	return response
}

func (i *MP3Library) GetAllRatings() map[string]int {
	ratings := map[string]int{}

	files, err := ioutil.ReadDir(i.BaseDir)
	if err != nil {
		log.Warningf("MP3Library.GetAllRatings ioutil.ReadDir: %v", err)
		return nil
	}

	for _, fs := range files {
		if fs.IsDir() {
			continue
		}
		if fs.Name() == "" {
			continue
		}

		fullPath, err := filepath.Abs(i.BaseDir + "/" + fs.Name())
		if err != nil {
			log.Warningf("MP3Library.GetAllRatings filepath.Abs: %v", err)
			return nil
		}

		rating := i.GetRating(fullPath)

		ratings[fs.Name()] = rating
	}

	return ratings
}
