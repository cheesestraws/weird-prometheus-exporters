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
	
	expected := map[string]SAS2IRCUAdapter{
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