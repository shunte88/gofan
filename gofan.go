package main

import (

    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
    "regexp"

    "github.com/TheSp1der/tplink"

    "periph.io/x/periph/conn/i2c/i2creg"
    "periph.io/x/periph/conn/physic"
    "periph.io/x/periph/devices/bmxx80"
    "periph.io/x/periph/host"

)

var (
    tphost  string
    address int
    lf      *os.File
    logf    *log.Logger
)

func init() {

    var err error
    lf, err = os.OpenFile("/tmp/gotemp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if nil!=err {
	log.Fatal(err)
    }
}

func main() {

    var trigger float64

    if nil!=lf {
        defer lf.Close()
        logf = log.New(lf, "gotemp::", log.LstdFlags)
    }
    re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

    trigp := flag.Float64("trigger", 30.00, "temperature trigger")
    flag.IntVar(&address, "address", 0x77, "address of temperature sensor")
    flag.StringVar(&tphost, "host", "", "hostname of device")
    flag.Parse()
    trigger = *trigp

    if len(strings.TrimSpace(tphost)) == 0 {
        fmt.Println("host must be defined")
        flag.PrintDefaults()
        os.Exit(1)
    }

    tp4mbh := tplink.Tplink{Host: tphost}
    status(tp4mbh)

    // Load all the drivers:
    if _, err := host.Init(); err != nil {
        log.Fatal(err)
        os.Exit(1) // never get here
    }

    // Open a handle to the first available I²C bus:
    bus, err := i2creg.Open("")
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    defer bus.Close()

    dev, err := bmxx80.NewI2C(bus, uint16(address), &bmxx80.DefaultOpts)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    defer dev.Halt()
    // Read temperature from the sensor:
    var env physic.Env
    if err = dev.Sense(&env); err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    str1 := fmt.Sprintf("%s", env.Temperature)
    cf := 0.00
    match := re.FindAllString(str1, -1)
    for _, ee := range match {
        fmt.Println(ee)
        cf, _ = strconv.ParseFloat(ee, 32)
        fmt.Println(cf)
    if cf > trigger {
        turnOn(tp4mbh)
    } else {
        turnOff(tp4mbh)
    }
    }

    s := fmt.Sprintf("%8s %10s %9s", env.Temperature, env.Pressure, env.Humidity)
    fmt.Println(s)
    logf.Println(s)
    status(tp4mbh)

}

func status(tp4mbh tplink.Tplink) {
    response, err := tp4mbh.SystemInfo()
    if err != nil {
        fmt.Println("unable to communicate with tp device")
        os.Exit(1)
    }
    s := ""
    if response.System.GetSysinfo.RelayState == 1 {
        s = "ON"
    } else {
        s = "OFF"
    }
    fmt.Println(tphost, "is", s)
    logf.Println(tphost, "is", s)
}

func turnOn(tp4mbh tplink.Tplink) {
    if err := tp4mbh.TurnOn(); err != nil {
        fmt.Println("unable to communicate with tp device")
        os.Exit(1)
    }
}

func turnOff(tp4mbh tplink.Tplink) {
    if err := tp4mbh.TurnOff(); err != nil {
        fmt.Println("unable to communicate with tp device")
        os.Exit(1)
    }
}

