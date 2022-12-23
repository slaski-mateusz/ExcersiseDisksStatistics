package main

import (
	"bytes"
	"reflect"
	"testing"
)

var dc_test_data = DisksData{
	DiskRecord{
		Datacenter:          "dc1",
		Hostname:            "host1",
		Serial:              "serial1",
		AgeSeconds:          1111111,
		TotalReads:          1111111,
		TotalWrites:         1111111,
		TotalReadsAndWrites: 1111111 + 1111111,
		AvgIoLatInMs:        1,
		TotalUncReadErr:     0,
		TotalUncWriteErr:    0,
	},
	DiskRecord{
		Datacenter:          "dc1",
		Hostname:            "host2",
		Serial:              "serial2",
		AgeSeconds:          2222222,
		TotalReads:          2222222,
		TotalWrites:         2222222,
		TotalReadsAndWrites: 2222222 + 2222222,
		AvgIoLatInMs:        2,
		TotalUncReadErr:     0,
		TotalUncWriteErr:    0,
	},
	DiskRecord{
		Datacenter:          "dc2",
		Hostname:            "host3",
		Serial:              "serial3",
		AgeSeconds:          3333333,
		TotalReads:          3333333,
		TotalWrites:         3333333,
		TotalReadsAndWrites: 3333333 + 3333333,
		AvgIoLatInMs:        3,
		TotalUncReadErr:     0,
		TotalUncWriteErr:    1,
	},
	DiskRecord{
		Datacenter:          "dc2",
		Hostname:            "host4",
		Serial:              "serial4",
		AgeSeconds:          4444444,
		TotalReads:          4444444,
		TotalWrites:         4444444,
		TotalReadsAndWrites: 4444444 + 4444444,
		AvgIoLatInMs:        20,
		TotalUncReadErr:     1,
		TotalUncWriteErr:    0,
	},
}

func TestContainsTrue(t *testing.T) {
	var slc []string
	var ele string
	// testing if element would be found
	slc = []string{"AAA", "BBB", "CCC"}
	ele = "CCC"
	if !Contains(slc, ele) {
		t.Fatalf("Contains not find element %v in slice %v", ele, slc)
	}
	// testing if not contained element can be found
	slc = []string{"AAA", "BBB", "CCC"}
	ele = "DDD"
	if Contains(slc, ele) {
		t.Fatalf("Contains found element %v in slice %v", ele, slc)
	}
}

func Test_read_file_data(t *testing.T) {
	var mock_buff bytes.Buffer
	mock_buff.WriteString("dc1;host1;sn1;1;2;3;4;5;6\ndc2;host2;sn2;6;5;4;3;2;1")
	wants := [][]string{
		[]string{
			"dc1",
			"host1",
			"sn1",
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
		},
		[]string{
			"dc2",
			"host2",
			"sn2",
			"6",
			"5",
			"4",
			"3",
			"2",
			"1",
		},
	}
	content := read_file_data(&mock_buff)
	if !reflect.DeepEqual(content, wants) {
		t.Fatalf("Reading CSV data unsucesfull")
	}
}

func Test_convert_line_to_disk_record(t *testing.T) {
	line := []string{
		"Data-center",
		"host-number.data-center.storage",
		"SERIAL",
		"1234567",
		"1234567",
		"1234567",
		"20",
		"0",
		"1",
	}
	want := DiskRecord{
		Datacenter:          "Data-center",
		Hostname:            "host-number.data-center.storage",
		Serial:              "SERIAL",
		AgeSeconds:          1234567,
		TotalReads:          1234567,
		TotalWrites:         1234567,
		TotalReadsAndWrites: 1234567 + 1234567,
		AvgIoLatInMs:        20,
		TotalUncReadErr:     0,
		TotalUncWriteErr:    1,
	}
	if convert_line_to_disk_record(line) != want {
		t.Fatalf("\nCSV line:\n%v\nnot parsed to expected disk record data:\n%v", line, want)
	}
}

func Test_count_disks_per_datacenter(t *testing.T) {
	var wants = "dc1: 2\ndc2: 2\n"
	if count_disks_per_datacenter(dc_test_data) != wants {
		t.Fatalf("Unsucessfull calculation of disks number per datacenter")
	}
}
