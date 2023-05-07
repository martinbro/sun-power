package sundata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

/* Soldata*/
type Hour struct {
	Time string  `josn:"time"`
	P    float64 `json:"P"`
}

type Hourly struct {
	Hourly []Hour `json:"hourly"`
}

type Pvgis struct {
	Outputs Hourly `json:"outputs"`
}

/* Kommunikationsdata */

type myConn struct {
	ch  chan []byte
	ch1 chan float64
	// msg chan[]byte
}

func SolData() {

	sol := myConn{ch: make(chan []byte), ch1: make(chan float64)} //, msg: make(chan []byte)}
	go prtSolData(sol)

	// s := fmt.Sprintf("https://re.jrc.ec.europa.eu/api/seriescalc?lat=54.854152&lon=10.504639&usehorizon=1&startyear=2016&endyear=2016&pvcalculation=1&peakpower=80&pvtechchoice=crystSi&mountingplace=building&loss=0.14&trackingtype=0&angle=0&aspect=0&outputformat=json") //DESC ASC
	s := "https://re.jrc.ec.europa.eu/api/seriescalc?lat=54.854152&lon=10.504639&usehorizon=1&startyear=2016&endyear=2016&pvcalculation=1&peakpower=80&pvtechchoice=crystSi&mountingplace=building&loss=0.14&trackingtype=0&angle=0&aspect=0&outputformat=json" //DESC ASC
	response, err := http.Get(s)
	check(err)

	jsonCunsumpData, err := io.ReadAll(response.Body)
	check(err)

	var pvgis Pvgis

	err = json.Unmarshal(jsonCunsumpData, &pvgis)
	check(err)

	layout := "20060102:1504"

	var max float64
	for i, p := range pvgis.Outputs.Hourly {
		if p.P > 0 {

			time, _ := time.Parse(layout, p.Time)
			//Konverterer til byte
			s := fmt.Sprintf("%.1f", p.P)
			b := []byte(s)
			if max < p.P {
				fmt.Println(i, ", String:", s, "Byte:", b, " Dato:", time)
				sol.ch <- b
				sol.ch1 <- p.P
			}
			max = math.Max(max, p.P)
		}
	}
}
func prtSolData(sol myConn) {
	for {
		a := <-sol.ch
		s := string(a)
		s1 := <-sol.ch1
		fmt.Println(a, len(a), s, len(s), s1)
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}

}
