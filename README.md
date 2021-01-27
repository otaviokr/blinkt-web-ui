# blinkt-web-ui

A web user interface to change the LEDs configuration in [Pimoroni's Blinkt!](https://shop.pimoroni.com/products/blinkt).

Blinkt! is an array of 8 individually addressable LEDs compatible with Raspberry Pi GPIO.

## Features

- Assign the color using native browser color picker for each one of the 8 LEDs in Blinkt array individually
- "Black" color is defined to turn the LED off
- Define global brightness of the LEDs
- No animation or any timely color changes

## How to run

`Keep in mind you will need to run this program with root privileges due to Blinkt library and Raspberry Pi GPIO limitations`

The easiest to get it running is to clone this repo and run the main source file. The web page will be available at [http://localhost:8090](http://localhost:8090).

```bash
git clone https://github.com/otaviokr/blinkt-web-ui.git
cd blinkt-web-ui
go run main.go
```

## How it works

You will be presented a page with 8 black circles, a color picker, a slider and a button:
- Each **circle** represents one of the LEDs in Blinkt. To change its color, select the color and click on the circle. The webpage has no feedback, so the color of the circle initially may differ from its current real state;
- Use the **color picker** to select the next color you want to configure the LEDs with. Black turns the LED off;
- The **slider** defines the brightness level of all LEDs. You cannot assign brightness individually to the LEDs, the last value before submitting the configuration will be the one used on Blinkt;
- After you've configured all the LEDs you want and set the desired brightness level, press the **button** to submit the values to Blinkt. The array should be updated shortly.

## Dependencies

All Blinkt handling is done via the excelllent Alex Ellis' [Blinkt lib in Go](https://github.com/alexellis/blinkt_go).

Logging is done with Sirupsen's [Logrus](https://github.com/sirupsen/logrus).

Look'n'Feel is powered by [Semantic-UI](https://semantic-ui.com/) and the [Superhero theme](https://github.com/semantic-ui-forest/forest-themes/blob/master/dist/bootswatch/v3/semantic.superhero.min.css). They are not required, of course, but nobody wants to use an ugly page...
