# resource-monitor
A tool to monitor any program's resource usage (memory, cpu)

## Usage
```
$ ./resource_monitor.sh -h
Usage: ./resource_monitor.sh [OPTIONS]

Options:
  --pid <PID>                   Monitor the process with the specified PID
  --program <PROGRAM_NAME>      Monitor the program with the specified name
  --sampling-interval <SECOND>  Set the sampling interval in seconds (default: 5)
  --output-file <FILE>          Specify the output file name (default: program_resource_usage.csv)
  -h, --help                    Show this help message
```

## Examples

Let's use etcd as an example to show how to monitor etcd's resource usage. See detailed steps below, or you can just
to execute the [example_monitor_etcd.sh](examples/etcd/example_monitor_etcd.sh) directly.

### Step 1: start etcd and play traffic

Firstly start etcd,
```
$ ./bin/etcd --quota-backend-bytes=8600000000
```

Afterwards, play traffic in another terminal,
```
$ bin/tools/benchmark txn-put --endpoints="http://127.0.0.1:2379" --clients=200 --conns=200 --key-space-size=4000000000 --key-size=128 --val-size=10240  --total=600000 --rate=40000
```

### Step 2: monitor the resource usage

When the etcd is running, execute any of the following command to monitor its resource usage,

- Monitor process with PID 70500

Assuming 70500 is the PID of etcd,

```
$./resource_monitor.sh --pid 70500
```

- Monitor program etcd
```
$./resource_monitor.sh --program etcd
```

- Monitor program etcd, and output the result into test_report_3.6.csv
```
$./resource_monitor.sh --program etcd --output-file test_report_3.6.csv
```

## Result
The result is in CSV format; an example is shown below.

```
Timestamp, Relative Time(sec), Program, PID, CPU(%), Memory(RSS KB)
2025-02-28 16:34:38, 0, ./bin/etcd, 4804, 0.0, 354944
2025-02-28 16:34:43, 5, ./bin/etcd, 4804, 0.1, 354944
2025-02-28 16:34:48, 10, ./bin/etcd, 4804, 0.0, 354944
2025-02-28 16:34:53, 15, ./bin/etcd, 4804, 0.0, 354944
2025-02-28 16:34:58, 20, ./bin/etcd, 4804, 222.3, 357376
2025-02-28 16:35:03, 25, ./bin/etcd, 4804, 234.2, 357456
```

You can render the data into line chart. See an example command below, you can input one or two csv files,

```
$ ./resource-monitor --data1 ./examples/etcd/testdata/test_report_3.5.18.csv --data2 ./examples/etcd/testdata/test_report_3.6.0-rc.1.csv
```

See the generated png file [line_chart.png](examples/etcd/testdata/line_chart.png)
