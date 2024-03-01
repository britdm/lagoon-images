if expr "$AIO_POD_METRICS_ENABLED" : '[Tt][Rr][Uu][Ee]' > /dev/null; then

    nohup /pod-metrics/pod-metrics-server &

fi