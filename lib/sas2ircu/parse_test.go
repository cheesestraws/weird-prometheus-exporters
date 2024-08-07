package sas2ircu

import (
	"testing"
	"reflect"
)

func TestParseSAS2IRCUList(t *testing.T) {
	output := `
LSI Corporation SAS2 IR Configuration Utility.
Version 19.00.00.00 (2014.03.17) 
Copyright (c) 2008-2014 LSI Corporation. All rights reserved. 


         Adapter      Vendor  Device                       SubSys  SubSys 
 Index    Type          ID      ID    Pci Address          Ven ID  Dev ID 
 -----  ------------  ------  ------  -----------------    ------  ------ 
   0     SAS2008     1000h    72h   00h:05h:00h:00h      1028h   1f1dh 
   1     SAS2009     1000h    72h   00h:06h:00h:00h      1029h   1f1dh
SAS2IRCU: Utility Completed Successfully.
	`
	
	stats := parseSAS2IRCUList([]byte(output))	
	
	expected := map[string]Adapter{
		"0": {
			Index: "0",
			AdapterType: "SAS2008",
			PCIAddress: "00h:05h:00h:00h",
		},
		"1": {
			Index: "1",
			AdapterType: "SAS2009",
			PCIAddress: "00h:06h:00h:00h",
		},
	}
	
	if !reflect.DeepEqual(expected, stats) {
		t.Errorf("didn't parse list properly")
	}
}

func TestIRStatusFromString(t *testing.T) {
	line := "floop (OKY) fleep"
	status := irStatusFromString(line)
	if status != IRStatusOkay {
		t.Errorf("where's my ok")
	}
	
	line = "floop (DGD) fleep"
	status = irStatusFromString(line)
	if status != IRStatusDegraded {
		t.Errorf("where's my degraded")
	}
}


func TestPDStatusFromString(t *testing.T) {
	line := "floop (OPT) fleep"
	status := pdStatusFromString(line)
	if status != PDStatusOptimal {
		t.Errorf("where's my ok")
	}
	
	line = "floop (DGD) fleep"
	status = pdStatusFromString(line)
	if status != PDStatusDegraded {
		t.Errorf("where's my degraded")
	}
}

func TestParseSAS2IRCUDisplay(t *testing.T) {
	output := `LSI Corporation SAS2 IR Configuration Utility.
Version 19.00.00.00 (2014.03.17) 
Copyright (c) 2008-2014 LSI Corporation. All rights reserved. 

Read configuration has been initiated for controller 0
------------------------------------------------------------------------
Controller information
------------------------------------------------------------------------
  Controller type                         : SAS2008
  BIOS version                            : 7.11.10.00
  Firmware version                        : 7.15.08.00
  Channel description                     : 1 Serial Attached SCSI
  Initiator ID                            : 0
  Maximum physical devices                : 39
  Concurrent commands supported           : 2607
  Slot                                    : 2
  Segment                                 : 0
  Bus                                     : 5
  Device                                  : 0
  Function                                : 0
  RAID Support                            : Yes
------------------------------------------------------------------------
IR Volume information
------------------------------------------------------------------------
IR volume 1
  Volume ID                               : 79
  Status of volume                        : Okay (OKY)
  Volume wwid                             : 00f6d172d82551b4
  RAID level                              : RAID1
  Size (in MB)                            : 953344
  Physical hard disks                     :
  PHY[0] Enclosure#/Slot#                 : 1:0
  PHY[1] Enclosure#/Slot#                 : 1:1
------------------------------------------------------------------------
Physical device information
------------------------------------------------------------------------
Initiator at ID #0

Device is a Hard disk
  Enclosure #                             : 1
  Slot #                                  : 0
  SAS Address                             : 4433221-1-0700-0000
  State                                   : Optimal (OPT)
  Size (in MB)/(in sectors)               : 953869/1953525167
  Manufacturer                            : ATA     
  Model Number                            : WDC WD1002FAEX-0
  Firmware Revision                       : 1D05
  Serial No                               : WDWCATRC693428
  GUID                                    : 50014ee2094077b8
  Protocol                                : SATA
  Drive Type                              : SATA_HDD

Device is a Hard disk
  Enclosure #                             : 1
  Slot #                                  : 1
  SAS Address                             : 4433221-1-0600-0000
  State                                   : Optimal (OPT)
  Size (in MB)/(in sectors)               : 953869/1953525167
  Manufacturer                            : ATA     
  Model Number                            : WDC WD1002FAEX-0
  Firmware Revision                       : 1D05
  Serial No                               : WDWMAY03993074
  GUID                                    : 50014ee002fe99b0
  Protocol                                : SATA
  Drive Type                              : SATA_HDD

Device is a Enclosure services device
  Enclosure #                             : 1
  Slot #                                  : 9
  SAS Address                             : 5e4ae02-0-b1db-7f00
  State                                   : Standby (SBY)
  Manufacturer                            : DP      
  Model Number                            : BACKPLANE       
  Firmware Revision                       : 1.07
  Serial No                               : 25500G7
  GUID                                    : N/A
  Protocol                                : SAS
  Device Type                             : Enclosure services device
------------------------------------------------------------------------
Enclosure information
------------------------------------------------------------------------
  Enclosure#                              : 1
  Logical ID                              : 5d4ae520:b1db7f00
  Numslots                                : 9
  StartSlot                               : 0
------------------------------------------------------------------------
SAS2IRCU: Command DISPLAY Completed Successfully.
SAS2IRCU: Utility Completed Successfully.
`
	devices := parseSAS2IRCUDisplay([]byte(output))
	
	expectedIRs := []IR{
		{
			VolumeID: "79",
			Status: IRStatusOkay,
		},
	}
	
	if !reflect.DeepEqual(devices.IRs, expectedIRs) {
		t.Errorf("IRs wrong: %v vs. %v", devices.IRs, expectedIRs)
	}
	
	expectedPhysDevs := []PhysicalDevice{
		{
			DeviceIsA: "Hard disk",
			Enclosure: "1",
			Slot: "0",
			State: PDStatusOptimal,
			SerialNumber: "WDWCATRC693428",
			Protocol: "SATA",
			DriveType: "SATA_HDD",
		},
		{
			DeviceIsA: "Hard disk",
			Enclosure: "1",
			Slot: "1",
			State: PDStatusOptimal,
			SerialNumber: "WDWMAY03993074",
			Protocol: "SATA",
			DriveType: "SATA_HDD",
		},
		{
			DeviceIsA: "Enclosure services device",
			Enclosure: "1",
			Slot: "9",
			State: PDStatusStandby,
			SerialNumber: "25500G7",
			Protocol: "SAS",
			DriveType: "",
		},
	}
	
	if !reflect.DeepEqual(devices.PhysicalDevices, expectedPhysDevs) {
		t.Errorf("phys devs wrong: %v vs. %v", devices.PhysicalDevices, expectedPhysDevs)
	}
}