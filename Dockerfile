# use prometheus as base layer
FROM prom/prometheus:v2.0.0
COPY bin/lnw /bin/lnw

ENV LNW_WATCH_PATH "/etc/prometheus/"

ENTRYPOINT [ "/bin/lnw", "/bin/prometheus" ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--web.console.libraries=/usr/share/prometheus/console_libraries", \
             "--web.console.templates=/usr/share/prometheus/consoles" ]
