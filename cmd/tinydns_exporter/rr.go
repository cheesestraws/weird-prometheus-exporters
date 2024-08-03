package main

import (
	"strconv"
	"strings"
)

type RRType int

const (
	RRTypeNone       RRType = 0
	RRTypeA          RRType = 1
	RRTypeNS         RRType = 2
	RRTypeMD         RRType = 3
	RRTypeMF         RRType = 4
	RRTypeCNAME      RRType = 5
	RRTypeSOA        RRType = 6
	RRTypeMB         RRType = 7
	RRTypeMG         RRType = 8
	RRTypeMR         RRType = 9
	RRTypeNULL       RRType = 10
	RRTypePTR        RRType = 12
	RRTypeHINFO      RRType = 13
	RRTypeMINFO      RRType = 14
	RRTypeMX         RRType = 15
	RRTypeTXT        RRType = 16
	RRTypeRP         RRType = 17
	RRTypeAFSDB      RRType = 18
	RRTypeX25        RRType = 19
	RRTypeISDN       RRType = 20
	RRTypeRT         RRType = 21
	RRTypeNSAPPTR    RRType = 23
	RRTypeSIG        RRType = 24
	RRTypeKEY        RRType = 25
	RRTypePX         RRType = 26
	RRTypeGPOS       RRType = 27
	RRTypeAAAA       RRType = 28
	RRTypeLOC        RRType = 29
	RRTypeNXT        RRType = 30
	RRTypeEID        RRType = 31
	RRTypeNIMLOC     RRType = 32
	RRTypeSRV        RRType = 33
	RRTypeATMA       RRType = 34
	RRTypeNAPTR      RRType = 35
	RRTypeKX         RRType = 36
	RRTypeCERT       RRType = 37
	RRTypeDNAME      RRType = 39
	RRTypeOPT        RRType = 41 // EDNS
	RRTypeAPL        RRType = 42
	RRTypeDS         RRType = 43
	RRTypeSSHFP      RRType = 44
	RRTypeIPSECKEY   RRType = 45
	RRTypeRRSIG      RRType = 46
	RRTypeNSEC       RRType = 47
	RRTypeDNSKEY     RRType = 48
	RRTypeDHCID      RRType = 49
	RRTypeNSEC3      RRType = 50
	RRTypeNSEC3PARAM RRType = 51
	RRTypeTLSA       RRType = 52
	RRTypeSMIMEA     RRType = 53
	RRTypeHIP        RRType = 55
	RRTypeNINFO      RRType = 56
	RRTypeRKEY       RRType = 57
	RRTypeTALINK     RRType = 58
	RRTypeCDS        RRType = 59
	RRTypeCDNSKEY    RRType = 60
	RRTypeOPENPGPKEY RRType = 61
	RRTypeCSYNC      RRType = 62
	RRTypeZONEMD     RRType = 63
	RRTypeSVCB       RRType = 64
	RRTypeHTTPS      RRType = 65
	RRTypeSPF        RRType = 99
	RRTypeUINFO      RRType = 100
	RRTypeUID        RRType = 101
	RRTypeGID        RRType = 102
	RRTypeUNSPEC     RRType = 103
	RRTypeNID        RRType = 104
	RRTypeL32        RRType = 105
	RRTypeL64        RRType = 106
	RRTypeLP         RRType = 107
	RRTypeEUI48      RRType = 108
	RRTypeEUI64      RRType = 109
	RRTypeNXNAME     RRType = 128
	RRTypeURI        RRType = 256
	RRTypeCAA        RRType = 257
	RRTypeAVC        RRType = 258
	RRTypeAMTRELAY   RRType = 260
)

var rrNames = map[RRType]string{
	0:   "None",
	1:   "A",
	2:   "NS",
	3:   "MD",
	4:   "MF",
	5:   "CNAME",
	6:   "SOA",
	7:   "MB",
	8:   "MG",
	9:   "MR",
	10:  "NULL",
	12:  "PTR",
	13:  "HINFO",
	14:  "MINFO",
	15:  "MX",
	16:  "TXT",
	17:  "RP",
	18:  "AFSDB",
	19:  "X25",
	20:  "ISDN",
	21:  "RT",
	23:  "NSAPPTR",
	24:  "SIG",
	25:  "KEY",
	26:  "PX",
	27:  "GPOS",
	28:  "AAAA",
	29:  "LOC",
	30:  "NXT",
	31:  "EID",
	32:  "NIMLOC",
	33:  "SRV",
	34:  "ATMA",
	35:  "NAPTR",
	36:  "KX",
	37:  "CERT",
	39:  "DNAME",
	41:  "OPT",
	42:  "APL",
	43:  "DS",
	44:  "SSHFP",
	45:  "IPSECKEY",
	46:  "RRSIG",
	47:  "NSEC",
	48:  "DNSKEY",
	49:  "DHCID",
	50:  "NSEC3",
	51:  "NSEC3PARAM",
	52:  "TLSA",
	53:  "SMIMEA",
	55:  "HIP",
	56:  "NINFO",
	57:  "RKEY",
	58:  "TALINK",
	59:  "CDS",
	60:  "CDNSKEY",
	61:  "OPENPGPKEY",
	62:  "CSYNC",
	63:  "ZONEMD",
	64:  "SVCB",
	65:  "HTTPS",
	99:  "SPF",
	100: "UINFO",
	101: "UID",
	102: "GID",
	103: "UNSPEC",
	104: "NID",
	105: "L32",
	106: "L64",
	107: "LP",
	108: "EUI48",
	109: "EUI64",
	128: "NXNAME",
	256: "URI",
	257: "CAA",
	258: "AVC",
	260: "AMTRELAY",
}

func (r RRType) String() string {
	return rrNames[r]
}

func TinyDNSLineToRRTypes(line string) []RRType {
	// source of truth: https://cr.yp.to/djbdns/tinydns-data.html

	if len(line) == 0 {
		return nil
	}

	switch line[0:1] {
	case ".":
		return []RRType{RRTypeNS, RRTypeSOA, RRTypeA}
	case "&":
		return []RRType{RRTypeNS, RRTypeA}
	case "=":
		return []RRType{RRTypeA, RRTypePTR}
	case "+":
		return []RRType{RRTypeA}
	case "@":
		return []RRType{RRTypeA, RRTypeMX}
	case "'":
		return []RRType{RRTypeTXT}
	case "^":
		return []RRType{RRTypePTR}
	case "C":
		return []RRType{RRTypeCNAME}
	case "Z":
		return []RRType{RRTypeSOA}
	case ":":
		// Arbitrary line.
		fields := strings.Split(line[1:], ":")
		if len(fields) < 2 {
			return nil
		}

		t, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil
		}

		return []RRType{RRType(t)}
	}
	return nil
}
