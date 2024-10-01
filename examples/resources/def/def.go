package def

var CAN_BE_ACCESS_FROM_PLUGIN = 10000

func init() { // will init run => yes
	CAN_BE_ACCESS_FROM_PLUGIN = -10000
}
