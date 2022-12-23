/*
The goal of this task is to prepare statistical analysis of set of data from disks.

Each entry of the data set consists of following fields separated by ;
character:

    datacenter
    hostname
    disk serial
    disk age (in s)
    total reads
    total writes
    average IO latency from 5 minutes (in ms)
    total uncorrected read errors
    total uncorrected write errors

The proper solution should output following information:

    How many disks are in total and in each DC
    Which disk is the youngest/oldest one and what is its age (in days)
    What's the average disk age per DC (in days)
    How many read/write IO/s disks processes on average
    Find top 5 disks with lowest/highest average IO/s (reads+writes, print disks and their avg IO/s)
    Find disks which are most probably broken, i.e. have non-zero uncorrected errors (print disks and error counter)

There should also be tests that verify if parts of the script are processing data properly.
*/

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

const (
	secondsInDay = 60 * 60 * 24
)

type DiskRecord struct {
	Datacenter          string
	Hostname            string
	Serial              string
	AgeSeconds          int
	TotalReads          int
	TotalWrites         int
	TotalReadsAndWrites int
	AvgIoLatInMs        int
	TotalUncReadErr     int
	TotalUncWriteErr    int
}

type DisksData []DiskRecord

func Contains[T comparable](sl []T, el T) bool {
	for _, val := range sl {
		if val == el {
			return true
		}
	}
	return false
}

func read_file_data(file_reader io.Reader) [][]string {
	csvReader := csv.NewReader(file_reader)
	csvReader.Comma = ';'
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Problem with reading CVS from file. Exiting!")
		os.Exit(2)
	}
	return data
}

func convert_line_to_disk_record(line []string) DiskRecord {
	age, _ := strconv.Atoi(line[3])
	totalReads, _ := strconv.Atoi(line[4])
	totalWrites, _ := strconv.Atoi(line[5])
	avg_io_lat_in_ms, _ := strconv.Atoi(line[6])
	totalUncReadErr, _ := strconv.Atoi(line[7])
	totalUncWriteErr, _ := strconv.Atoi(line[8])
	return DiskRecord{
		Datacenter:          line[0],
		Hostname:            line[1],
		Serial:              line[2],
		AgeSeconds:          age,
		TotalReads:          totalReads,
		TotalWrites:         totalWrites,
		TotalReadsAndWrites: totalReads + totalWrites,
		AvgIoLatInMs:        avg_io_lat_in_ms,
		TotalUncReadErr:     totalUncReadErr,
		TotalUncWriteErr:    totalUncWriteErr,
	}
}

func convert_raw_to_disk_records(in_data [][]string) DisksData {
	var outData DisksData
	for _, line := range in_data {
		outData = append(outData, convert_line_to_disk_record(line))
	}
	return outData
}

func count_disks_per_datacenter(in_data DisksData) string {
	disksPerDc := make(map[string]int)
	for _, dr := range in_data {
		if _, ok := disksPerDc[dr.Datacenter]; !ok {
			disksPerDc[dr.Datacenter] = 0
		}
		disksPerDc[dr.Datacenter] = disksPerDc[dr.Datacenter] + 1
	}
	out, _ := yaml.Marshal(disksPerDc)
	return string(out)
}

func average_age_per_datacenter(in_data DisksData) string {
	sumForDc := make(map[string]int)
	countForDc := make(map[string]int)
	averageForDc := make(map[string]float32)
	for _, dr := range in_data {
		sumForDc[dr.Datacenter] += dr.AgeSeconds
		countForDc[dr.Datacenter]++
	}
	for dcn, sum := range sumForDc {
		averageForDc[dcn] = (float32(sum) / float32(countForDc[dcn])) / float32(secondsInDay)
	}
	out, _ := yaml.Marshal(averageForDc)
	return string(out)
}

func count_average_io(in_data DisksData) string {
	var readsSum int = 0
	var writesSum int = 0
	var count int = 0
	for _, dr := range in_data {
		readsSum += dr.TotalReads
		writesSum += dr.TotalWrites
		count++
	}
	out, _ := yaml.Marshal(map[string]float32{
		"average_reads":  float32(readsSum) / float32(count),
		"average_writes": float32(writesSum) / float32(count),
	})
	return string(out)
}

type OldYoungDisks struct {
	Oldest   DiskRecord
	Youngest DiskRecord
}

func find_oldest_and_youngest_disks(in_data DisksData) string {
	dStats := new(OldYoungDisks)
	for idx, dr := range in_data {
		if idx == 1 {
			dStats.Oldest = dr
			dStats.Youngest = dr
			continue
		}
		if dr.AgeSeconds < dStats.Youngest.AgeSeconds {
			dStats.Youngest = dr
		}
		if dr.AgeSeconds > dStats.Oldest.AgeSeconds {
			dStats.Oldest = dr
		}
	}
	out, _ := yaml.Marshal(*dStats)
	return string(out)
}

func (dire DiskRecord) age_days() float32 {
	return float32(dire.AgeSeconds) / float32(secondsInDay)
}

func rank_read_write_io(in_data DisksData, ndisks int) string {
	var mostLoaded DisksData
	var leastLoaded DisksData
	for cyc := 0; cyc < ndisks; cyc++ {
		var cyc_max DiskRecord
		var cyc_min DiskRecord
		for _, dr := range in_data {
			if !Contains(mostLoaded, dr) {
				if (cyc_max == DiskRecord{}) {
					cyc_max = dr
				}
				if cyc_max.TotalReadsAndWrites < dr.TotalReadsAndWrites {
					cyc_max = dr
				}
			}
			if !Contains(leastLoaded, dr) {
				if (cyc_min == DiskRecord{}) {
					cyc_min = dr
				}
				if cyc_min.TotalReadsAndWrites > dr.TotalReadsAndWrites {
					cyc_min = dr
				}
			}
		}
		mostLoaded = append(mostLoaded, cyc_max)
		leastLoaded = append(leastLoaded, cyc_min)
	}
	out, _ := yaml.Marshal(map[string]DisksData{
		"Most loaded":  mostLoaded,
		"Least loaded": leastLoaded,
	})
	return string(out)
}

func find_broken_disks(in_data DisksData) string {
	var brdi DisksData
	for _, dr := range in_data {
		if dr.TotalUncReadErr > 0 || dr.TotalUncWriteErr > 0 {
			brdi = append(brdi, dr)
		}
	}
	out, _ := yaml.Marshal(brdi)
	return string(out)
}

func main() {
	fileName := flag.String("filename", "data.raw", "Input CSV filename")
	flag.Parse()
	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("Can`t open file %v\n", *fileName)
		os.Exit(1)
	}
	disks_raw_data := read_file_data(file)
	disks_data := convert_raw_to_disk_records(disks_raw_data)
	fmt.Println("Number of disks per datacenter:")
	fmt.Println(count_disks_per_datacenter(disks_data))
	fmt.Println("Oldest and youngest disks:")
	fmt.Println(find_oldest_and_youngest_disks(disks_data))
	fmt.Println("Average disk age in days per datacenter")
	fmt.Println(average_age_per_datacenter(disks_data))
	fmt.Println("Average IO")
	fmt.Println(count_average_io(disks_data))
	fmt.Println("Top loaded and lazy disks")
	fmt.Println(rank_read_write_io(disks_data, 5))
	fmt.Println("Broken disks")
	fmt.Println(
		find_broken_disks(disks_data),
	)
}
