#!/usr/bin/env bash

cleanup() {
    kill $MONITOR_PID
    kill $BENCHMARK_PID
    kill $ETCD_PID
}

trap cleanup SIGINT

# start etcd
./etcd --quota-backend-bytes=8600000000 > etcd_output.log 2>&1 &
ETCD_PID=$!

sleep 2

# play traffic
./benchmark txn-put --endpoints="http://127.0.0.1:2379" --clients=200 --conns=200 --key-space-size=4000000000 --key-size=128 --val-size=10240  --total=600000 --rate=40000 > benchmark_output.log 2>&1 &
BENCHMARK_PID=$!

sleep 2

# monitor etcd's resource usage
./resource_monitor.sh --program etcd &
MONITOR_PID=$!

# wait for the benchmark test to finish
wait $BENCHMARK_PID

kill $MONITOR_PID
kill $ETCD_PID

echo "Done"
