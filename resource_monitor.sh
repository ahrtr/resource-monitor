#!/usr/bin/env bash

# Copyright (c) 2025 Benjamin Wang
# Licensed under the MIT License

# Examples
#
# Show help
#   $./resource_monitor.sh -h
#
# Monitor process with PID 70500
#   $./resource_monitor.sh --pid 70500
#
# Monitor program etcd
#   $./resource_monitor.sh --program etcd
#
# Monitor program etcd, and output the result into test_report.csv
#   $./resource_monitor.sh --program etcd --output-file test_report.csv

# Default values
default_sampling_interval=5
pid=""
program_name=""
sampling_interval=""
output_file="program_resource_usage.csv"

show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --pid <PID>                   Monitor the process with the specified PID"
    echo "  --program <PROGRAM_NAME>      Monitor the program with the specified name"
    echo "  --sampling-interval <SECOND>  Set the sampling interval in seconds (default: $default_sampling_interval)"
    echo "  --output-file <FILE>          Specify the output file name (default: $output_file)"
    echo "  -h, --help                    Show this help message"
    echo ""
}

# parse command line arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            show_help
            exit 0
            ;;
        --pid)
            pid="$2"
            shift 2
            ;;
        --program)
            program_name="$2"
            shift 2
            ;;
        --sampling-interval)
            sampling_interval="$2"
            shift 2
            ;;
        --output-file)
            output_file="$2"
            shift 2
            ;;    
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

if [ -z "$pid" ] && [ -z "$program_name" ]; then
    echo "Error: You must specify either --pid or --program."
    show_help
    exit 1
fi

# Use the default value if not set in command line parameter
sampling_interval="${sampling_interval:-$default_sampling_interval}"

# Create the output file and output the head line if not present
if [ ! -f "$output_file" ]; then
    echo "Timestamp, Relative Time(sec), Program, PID, CPU(%), Memory(RSS KB)" > "$output_file"
fi

# record the timestamp of the first sampling, which is used to calculate the relative time
start_time=$(date +%s)  

sample_count=0
while true
do
    # Get current time
    timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    current_time=$(date +%s)
    relative_time=$((current_time - start_time))

    # Monitor the process specified by $pid if set
    if [ ! -z "$pid" ]; then
        # exclude the headline
        ps_output=$(ps -p "$pid" -o comm,pid,%cpu,rss | tail -n +2)
    elif [ ! -z "$program_name" ]; then
        # filter by the program name
        ps_output=$(ps -eo comm,pid,%cpu,rss | grep "$program_name" | grep -v "grep")
    fi

    if [ ! -z "$ps_output" ]; then
        # convert into CSV format
        echo "$timestamp, $relative_time, $(echo "$ps_output" | awk '{print $1", "$2", "$3", "$4}')" >> "$output_file"
    else
        # Stop monitoring if the program isn't running
        if [ ! -z "$pid" ]; then
            echo "Program with PID $pid is not running. Stopping monitoring."
        elif [ ! -z "$program_name" ]; then
            echo "Program $program_name is not running. Stopping monitoring."
        else
            echo "No matching programs found. Stopping monitoring."
        fi
        exit 0
    fi

    sleep $sampling_interval
    sample_count=$((sample_count + 1))
done
