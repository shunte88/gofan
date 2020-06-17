# Go Fan Go!

_Copyright &copy; Stuart Hunter. All rights reserved._

Temperature trigger that runs on a Raspberry Pi Zero W

Simple automation mechanism to trigger a TP-Link Kasa Smart Plug over wi-fi when a specified temperature limit is reached

The Pi uses an I2C BMP180, BMP280 superceeds, to determine the temperature.

* Command line specification of trigger temperature and Kasa address
* Trigger logic can be [A]bove or [B]elow specified limit
* Scheduled to run at intervals, cron keeps it simple

## Installation

### Installing Dependencies

1. Enable I2C via raspi-config
2. Install latest supported go or tiny go
3. Install periph.io
4. Install tplink library

### Installing GoFan

```sh
go get -d github.com/TheSp1der/tplink
git clone https://github.com/shunte88/gofan.git
cd gofan
go build
```

## Configuration

Configure GoFan to run at intervals using cron, this provides a simple and robust solution

```sh
*/20 * * * * /home/pi/gofan/gofan -host tp_br1 -trigger 28.00 --logic A
```

Note that entries in /etc/hosts have been added here to echo the names given to each Kasa device in the Kasa mobile app, alternatively the IP address of the device may be specified

## TODO


