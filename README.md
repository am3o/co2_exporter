# airco2ntrol CO₂ Exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/am3o/co2_exporter)](https://goreportcard.com/report/github.com/am3o/co2_exporter)
[![Docker Image](https://img.shields.io/badge/ghcr.io-am3o%2Fco2--exporter-blue?logo=github)](https://ghcr.io/am3o/co2_exporter)
[![Release](https://github.com/am3o/co2_exporter/actions/workflows/release.yml/badge.svg)](https://github.com/am3o/co2_exporter/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Prometheus exporter for collecting CO₂, temperature, and humidity metrics from the  
[AIRCO2NTROL COACH](https://www.tfa-dostmann.de/produkt/co2-monitor-airco2ntrol-coach-31-5009/) USB monitor.

This exporter communicates with the device via the Linux HIDRAW interface (e.g., `/dev/hidraw0`).

---

## 📈 Exported Metrics

| Metric Name                            | Type  | Labels | Unit     |
|----------------------------------------|:-----:|:------:|----------|
| `tfa_airco2ntrol_carbon_dioxide_total` | Gauge | `unit` | ppm      |
| `tfa_airco2ntrol_temperature_total`    | Gauge | `unit` | °C       |
| `tfa_airco2ntrol_humidity_total`       | Gauge | `unit` | percent  |

---

## 🚀 Installation & Usage

### From Source

```bash
go install github.com/am3o/co2_exporter@latest

# Run with root privileges due to HIDRAW access requirement
sudo CO2MOINITOR_DEVICE=/dev/hidraw0 $GOPATH/bin/co2_exporter
```
> [!NOTE]
> Accessing HIDRAW devices typically requires root permissions.

### Using Docker
```bash
docker run -d \
    --name co2-exporter \
    -v /dev/hidraw0:/dev/hidraw0:ro \
    --privileged \
    -p 8080:8080 \
    ghcr.io/am3o/co2_exporter:latest
```

## 📡 Prometheus Configuration

To scrape metrics from co2_exporter, add the following to your Prometheus configuration (`prometheus.yml`):
```yaml
scrape_configs:
  - job_name: 'co2_exporter'
    static_configs:
      - targets: ['localhost:8080']
```

Adjust the target if you're running the exporter on a different host or port.

## 🛠️ systemd Service Template

Save this file as `/etc/systemd/system/co2_exporter.service`:
```ini
[Unit]
Description=CO₂ Exporter
After=network.target

[Service]
ExecStart=/usr/local/bin/co2_exporter
Environment=CO2MOINITOR_DEVICE=/dev/hidraw0
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target
```

Enable and start the service: 

```bash
sudo systemctl daemon-reexec
sudo systemctl enable --now co2_exporter.service
```
## 🤝 Contributing

Contributions are welcome! Feel free to fork the repository and submit a pull request.
