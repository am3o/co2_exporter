# co2 exporter

The exporter provides functionality to export values from USB CO₂ monitor:

* [AIRCO2NTROL COACH](https://www.tfa-dostmann.de/produkt/co2-monitor-airco2ntrol-coach-31-5009/)

The core of the service use the API from the system Linux HIDRAW (`/dev/hidraw0` etc.) to access the USB CO₂ monitor.

## Metrics

| Name                                 | Type  | Label |    Unit |
|--------------------------------------|:-----:|:-----:|--------:|
| tfa_airco2ntrol_carbon_dioxide_total | Gauge | Unit  |     ppm |
| tfa_airco2ntrol_temperature_total    | Gauge | Unit  | celsius |
| tfa_airco2ntrol_humidity_total       | Gauge | Unit  | percent |

## Installation and usage

### From source

```bash
go get github.com/netzaffe/co2_exporter

$ sudo CO2MOINITOR_DEVICE=/dev/hidraw0 $GOPATH/bin/co2_exporter
```

Hint: sadly, to access the HIDRAW API you need root permissions.

### docker

```bash
docker run -d --name co2-exporter -v /dev/hidraw0:/dev/hidraw0:ro --privileged -p 8080:8080 netzaffe/co2-exporter:latest
```

## Contribute

Feel free to contribute to project.