package protocol

// Maximum supported Quality of Service
const (
	MAXQoS byte = 0x1
)

// Control packet codes ( shifted to left, mask : 0xF0 )
const (
	CCONNECT     byte = byte(0x1 << 4)
	CCONNACK     byte = byte(0x2 << 4)
	CQUEUE       byte = byte(0x4 << 4)
	CQUEUEACK    byte = byte(0x5 << 4)
	CPUBACK      byte = byte(0x6 << 4)
	CSUBSCRIBE   byte = byte(0x7 << 4)
	CSUBACK      byte = byte(0x8 << 4)
	CUNSUBSCRIBE byte = byte(0x9 << 4)
	CUNSUBACK    byte = byte(0xA << 4)
	CPUBLISH     byte = byte(0xB << 4)
	CPING        byte = byte(0xC << 4)
	CPONG        byte = byte(0xD << 4)
	CDISCONNECT  byte = byte(0xE << 4)
	// TODO
	//  CRESACK byte      = byte(0x6 << 4)
	//  CREQACK      byte = byte(0x4 << 4)
	//  CREQUEST     byte = byte(0x5 << 4)
	// 	CRESPONSE    byte = byte(0x3 << 4)
)

// Quality of Service codes
const (
	LQOS0 = 0x00
	LQOS1 = 0x01
	LQOS2 = 0x02
)

// Duplicate option
const (
	NDUP = 0
	YDUP = 1
)

// Retain option ( N = no, Y = yes . ex. NRET = no retain, YRET = retain)
const (
	NRET = 0
	YRET = 1
)

// Connection response header options
const (
	// TODO
	// . add the rest
	RHASSESSION = 0x08
	RCLEANSTART = 0x04
)

// Control packet raw codes
const (
	PNULL        = 0x00
	PCONNECT     = 0x01
	PCONNACK     = 0x02
	PQUEUE       = 0x04
	PQUEUEACK    = 0x05
	PPUBACK      = 0x06
	PSUBSCRIBE   = 0x07
	PSUBACK      = 0x08
	PUNSUBSCRIBE = 0x09
	PUNSUBACK    = 0x0A
	PPUBLISH     = 0x0B
	PPING        = 0x0C
	PPONG        = 0x0D
	PDISCONNECT  = 0x0E
	// TODO
	//  PRESACK      = 0x06
	// NOTE: new control codes should be included
	//
	// RREQUEST
	// RRESPONSE
	// RREQBCST
	// REQRBCST
	// PPROPOS
	// PRPROPS
	// PRESPONSE    = 0x03
	// PREQACK      = 0x04
	// PREQUEST     = 0x05
)
