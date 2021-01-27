package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"

	. "github.com/alexellis/blinkt_go"
)

var (
	blinkt Blinkt
)

// ParseColor will convert the hexadecimal in the string into a valid int.
func ParseColor(c string) int {
	r, err := strconv.ParseInt(c, 16, 32)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ColorComponent": c,
			},
		).Error("Could not parse the color component")
		return 0
	}
	return int(r)
}

// SetPixel will parse the incomming details for a LED (a "pixel") to configure it correctly.
func SetPixel(led, rgb string) (int, int, int, int) {
	index, err := strconv.Atoi(led[len(led) - 1:])
	if err != nil {
		log.WithFields(
			log.Fields{
				"led": led,
				"rgb": rgb,
			},
		).Error("Could not convert LED index")
		return -1, 0, 0, 0
	}

	re := regexp.MustCompile("#([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})")
	match := re.FindStringSubmatch(rgb)

	red := 0
	green := 0
	blue := 0
	if len(match) == 4 {
		red = ParseColor(match[1])
		green = ParseColor(match[2])
		blue = ParseColor(match[3])
	} else {
		log.WithFields(
			log.Fields{
				"led": led,
				"rgb": rgb,
			},
		).Error("Could not parse RGB")
	}

	return index - 1, red, green, blue
}

// UpdateLed will collect the LED values and brightness and send them to Blinkt.
func UpdateLed(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
			panic(err)
	}

	for name, value := range req.PostForm {
		if name == "bright" {
			br, err := strconv.Atoi(value[0])
			if err != nil {
				panic(err)
			}
			brightness := float64(br)/100
			log.WithFields(
				log.Fields{
					"brightness": brightness,
				},
			).Info("Set new value to brightness")
			blinkt.SetBrightness(brightness)
		} else {
			i, r, g, b := SetPixel(name, value[0])
			log.WithFields(
				log.Fields{
					"LedIndex": i,
					"Red": r,
					"Green": g,
					"Blue": b,
				},
			).Info("Set new color to LED")
			blinkt.SetPixel(i, r, g, b)
		}
	}

	blinkt.Show()
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Starting server")

	brightness := 0.125
	log.WithFields(
		log.Fields{
			"brightness": brightness,
		},
	).Info("Set default brightness")

	blinkt = NewBlinkt(brightness)
	blinkt.Setup()
	blinkt.SetClearOnExit(true)

	http.HandleFunc("/update_led", UpdateLed)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	port := 8090
	log.WithFields(
		log.Fields{
			"port": port,
		},
	).Info("Starting to listen at port")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
