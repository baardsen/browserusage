package dao

import (
	"code.google.com/p/goauth2/oauth/jwt"
	"code.google.com/p/google-api-go-client/analytics/v3"
	"encoding/gob"
	"encoding/json"
	resourcelocator "github.com/baardsen/resourcelocator"
	"log"
	"os"
	"strconv"
	"time"
)

type Point struct {
	Time  int64 `json:"x"`
	Value int   `json:"y"`
}

type Series struct {
	Name string  `json:"name"`
	Data []Point `json:"data"`
}

type measure struct {
	Name  string
	Time  int64
	Value int
}

var config struct {
	Profile         string
	ClientId        string
	ClientEmail     string
	CertificateFile string
	Certificate     []byte
	TokenUri        string
}

func Init() {
	tokData := resourcelocator.Locate("/resources/config.json")
	if err := json.Unmarshal(tokData, &config); err != nil {
		log.Fatal("dao.init", err)
	}
	config.Certificate = resourcelocator.Locate(config.CertificateFile)

	cacheFile := os.TempDir() + string(os.PathSeparator) + "browserusage.dat"
	if file, err := os.Open(cacheFile); err != nil {
		log.Println("Couldn't open cacheFile: "+cacheFile, err)
	} else {
		decoder := gob.NewDecoder(file)
		decoder.Decode(&cache)
		file.Close()
	}
	Query(firstDate, time.Now())

	if file, err := os.Create(cacheFile); err != nil {
		log.Println("Couldn't create cacheFile: "+cacheFile, err)
	} else {
		encoder := gob.NewEncoder(file)
		encoder.Encode(cache)
		file.Close()
	}
}

var firstDate = time.Date(2012, 4, 16, 0, 0, 0, 0, time.UTC)

func Query(from, to time.Time) []Series {
	from = from.In(time.UTC).AddDate(0, 0, int(time.Monday-from.Weekday()))
	if from.Before(firstDate) {
		from = firstDate
	}
	if to.IsZero() || to.After(time.Now()) {
		to = time.Now()
	}
	dataGaService := createDataService()
	ch := make(chan measure)
	done := make(chan struct{})
	counter := 0
	for from.Before(to) {
		counter++
		go makeRequest(dataGaService, from, ch, done)
		if cache[from] == nil {
			time.Sleep(1 * time.Second)
		}
		from = from.AddDate(0, 0, 7)
	}
	values := make(map[string][]Point, 0)
	for {
		select {
		case m := <-ch:
			browser := m.Name
			points := values[browser]
			if points == nil {
				points = make([]Point, 0)
			}
			values[browser] = append(points, Point{m.Time, m.Value})
		case <-done:
			counter--
			if counter == 0 {
				close(ch)
				close(done)
				series := make([]Series, 0)
				for name, points := range values {
					series = append(series, Series{name, points})
				}
				return sortSeries(series)
			}
		}
	}
}

const format = "2006-01-02"

var cache = make(map[time.Time][][]string)

func makeRequest(dataGaService *analytics.DataGaService, from time.Time, ch chan measure, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()
	if from.AddDate(0, 0, 7).After(time.Now()) {
		return
	}
	rows := cache[from]
	if rows == nil {
		log.Println("Sending request: " + from.Format(format))
		dataGaGetCall := dataGaService.Get(config.Profile, from.Format(format), from.AddDate(0, 0, 6).Format(format), "ga:users")
		dataGaGetCall.Dimensions("ga:browser")
		gaData, err := dataGaGetCall.Do()
		if err != nil {
			log.Fatal("GaData:", err)
		}
		rows = gaData.Rows
		cache[from] = rows
	}
	for _, arr := range rows {
		count, _ := strconv.Atoi(arr[1])
		ch <- measure{arr[0], from.Unix() * 1000, count}
	}
}

func createDataService() *analytics.DataGaService {
	// Craft the ClaimSet and JWT token.
	token := jwt.NewToken(config.ClientEmail, analytics.AnalyticsReadonlyScope, config.Certificate)
	//token.ClaimSet.Aud = config.TokenUri

	transport, err := jwt.NewTransport(token)
	if err != nil {
		log.Fatalf("failed to create authenticated transport: %+v.", err)
	}
	analyticsService, _ := analytics.New(transport.Client())
	return analytics.NewDataGaService(analyticsService)
}
