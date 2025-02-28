# resource-monitor
A tool to monitor any program's resource usage (memory, cpu)

## Usage
```
F5WPP69Q7H:resource-monitor wachao$ ./resource_monitor.sh -h
Usage: ./resource_monitor.sh [OPTIONS]

Options:
  --pid <PID>                   Monitor the process with the specified PID
  --program <PROGRAM_NAME>      Monitor the program with the specified name
  --sampling-interval <SECOND>  Set the sampling interval in seconds (default: 5)
  --output-file <FILE>          Specify the output file name (default: program_resource_usage.csv)
  -h, --help                    Show this help message
```

## Examples

- Monitor process with PID 70500
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
