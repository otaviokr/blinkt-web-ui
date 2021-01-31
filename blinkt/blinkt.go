package blinkt

import (
	"fmt"
	"log"
	"periph.io/x/conn/gpio"
	"periph.io/x/host"
	"periph.io/x/conn/gpio/gpioreg"
	"os"
	"os/signal"
	"time"
)

const (
	DataPin  string = "23" // "GPIO23"
	ClockPin string = "24" // "GPIO24"
)

// Pixel represents a LED in Blinkt (which has a total of 8 LEDs)
type Pixel struct {
	R          int
	G          int
	B          int
	Brightness float64
}

// Dev represents a Blinkt device (https://shop.pimoroni.com/products/blinkt)
type Dev struct {
	Dat   gpio.PinOut
	Clk   gpio.PinOut
	Array []Pixel
}

// Init is just a wrapper for host.Init()
func Init() error {
	_, err := host.Init()
	return err
}

// NewDev creates a new instance of Dev
func NewDev() (*Dev, error) {
	dat := gpioreg.ByName(DataPin)
	if dat == nil {
		log.Fatalf("Could not assign DAT to pin %s", DataPin)
	}

	clk := gpioreg.ByName(ClockPin)
	if clk == nil {
		log.Fatalf("Could not assign CLK to pin %s", ClockPin)
	}

	dev := &Dev{
		Dat:   dat,
		Clk:   clk,
		Array: make([]Pixel, 8),
	}

	return dev, nil
}

// SetClearOnExit turns all pixels off with Ctrl+C and/or os.Interrupt signal
func (d *Dev) SetClearOnExit(clearOnExit bool) {
	if clearOnExit {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		fmt.Println("Press Ctrl+C to stop")

		go func() {
			for range signalChan {
				d.Clear()
				d.Show()
				os.Exit(1)
			}
		}()
	}
}

// Delay is just a wrapper for time.Sleep
func Delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// SetPixel configures the LED at position i with color (r,g,b) and brightness
func (d *Dev) SetPixelWithBright(i, r, g, b int, brightness float64) {
	d.Array[i] = Pixel{
		R:          r,
		G:          g,
		B:          b,
		Brightness: brightness,
	}
}

func (d *Dev) SetPixelColor(i, r, g, b int) {
	d.Array[i].R = r
	d.Array[i].G = g
	d.Array[i].B = b
}

func (d *Dev) SetPixelBrightness(i int, brightness float64) {
	d.Array[i].Brightness = brightness
}

// Show sends all LEDs configuration to Blinkt device to update the LEDs
// We need to send 32 pulses with DAT=0 before sending the LED configuration
// For each LED, we need to send 3 bits set (111), then 5 bits for brightness
// Then 8 bits for Blue, 8 bits for Green and 8 bits for Red
// Finally, we send 36 pulses with DAT=0
func (d *Dev) Show() {
	// Send 32 pulses in Clock before start sending data.
	d.Dat.Out(gpio.Low)
	for i := 0; i < 32; i++ {
		d.Clk.Out(gpio.High)
		d.Clk.Out(gpio.Low)
	}

	// Send RGB to LEDs.
	for _, p := range d.Array {
		bitwise := 224 // 0b11100000
		d.Write(int(p.Brightness*31.0) | bitwise)
		d.Write(p.B)
		d.Write(p.G)
		d.Write(p.R)
	}

	// Sends 36 pulses in Clock to finish transmission.
	d.Dat.Out(gpio.Low)
	for i := 0; i < 36; i++ {
		d.Clk.Out(gpio.High)
		d.Clk.Out(gpio.Low)
	}
}

// Write sends the correct level to each GPIO (CLK/DAT) according to the value passed
func (d *Dev) Write(value int) {
	for i := 0; i < 8; i++ {
		// 0b10000000 = 128
		// 128 is a filter because we are only interested in the leftmost bit at this moment
		if value&128 == 0 {
			d.Dat.Out(gpio.Low)
		} else {
			d.Dat.Out(gpio.High)
		}
		d.Clk.Out(gpio.High)

		// Shift the value (second leftmost bit becomes the first leftmost bit)
		value = value << 1

		d.Clk.Out(gpio.Low)
	}
}

// SetAllPixels is juat a wrapper to call SetPixel() for all pixels, to set them equally
func (d *Dev) SetAllPixels(r, g, b int, brightness float64) {
	for i, _ := range d.Array {
		d.Array[i] = Pixel {
			R: r,
			G: g,
			B: b,
			Brightness: brightness,
		}
	}
}

// Clear turns off all LEDs (set all LEDs to (0, 0, 0))
func (d *Dev) Clear() {
	d.SetAllPixels(0, 0, 0, 0.0)
	d.Show()
}
