package main

import (
	"fmt"
	"github.com/otaviokr/blinkt-web-ui/blinkt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
)

var (
	blinktDev *blinkt.Dev
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
	index, err := strconv.Atoi(led[len(led)-1:])
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

	p := []blinkt.Pixel{
		blinkt.Pixel{}, blinkt.Pixel{}, blinkt.Pixel{}, blinkt.Pixel{},
		blinkt.Pixel{}, blinkt.Pixel{}, blinkt.Pixel{}, blinkt.Pixel{}}

	for name, value := range req.PostForm {
		if name[:5] == "input" {
			// Field with name ending with "b" contains brights for that led
			if name[len(name)-1:] == "b" {
				br, err := strconv.Atoi(value[0])
				if err != nil {
					panic(err)
				}
				brightness := float64(br) / 100
				led, err := strconv.Atoi(name[7:len(name)-1])
				led -= 1
				if err != nil {
					panic(err)
				}
				log.WithFields(
					log.Fields{
						"led":        led,
						"brightness": brightness,
					},
				).Info("Set new value to brightness to LED")
				p[led].Brightness = brightness
			} else {
				i, r, g, b := SetPixel(name, value[0])
				log.WithFields(
					log.Fields{
						"LedIndex": i,
						"Red":      r,
						"Green":    g,
						"Blue":     b,
					},
				).Info("Set new color to LED")
				p[i].R = r
				p[i].G = g
				p[i].B = b
			}
		}
	}

	for pi, pd := range p {
		blinktDev.SetPixelWithBright(pi, pd.R, pd.G, pd.B, pd.Brightness)
	}
	blinktDev.Show()
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Starting server")

	err := blinkt.Init()
	if err != nil {
		panic(err)
	}

	blinktDev, err = blinkt.NewDev()
	if err != nil {
		panic(err)
	}
	blinktDev.SetClearOnExit(true)

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
