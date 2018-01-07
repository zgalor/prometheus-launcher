FROM prom/prometheus:v2.0.0

COPY launcher /bin/launcher

ENTRYPOINT [ "/bin/launcher", "--watch-path=/etc/prometheus/", "/bin/prometheus"  ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml" ]
