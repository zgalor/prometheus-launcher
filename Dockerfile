# use prometheus as base layer
FROM prom/prometheus:v2.0.0
COPY bin/lnw /bin/lnw

ENTRYPOINT [ "/bin/lnw", "--watch-path=/etc/prometheus/", "/bin/prometheus"  ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--web.console.libraries=/usr/share/prometheus/console_libraries", \
             "--web.console.templates=/usr/share/prometheus/consoles" ]
