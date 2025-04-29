FROM scratch

COPY co2_exporter /co2_exporter

ENTRYPOINT ["/co2_exporter"]