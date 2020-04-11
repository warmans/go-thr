package amp

// Command that will tell the amp to send events back to the host whenever a control is changed on the
// physical amp.
var EnableEvents = CommandSet{
		{Type: msgTypeOne, Payload: []byte{0x07, 0x00, 0x04, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Type: msgTypeOne, Payload: []byte{0x03, 0x28, 0x24, 0x6b, 0x09, 0x18, 0x00, 00, 0x00}},
}

// GetActivePreset will return what the currently active channel.
var GetActivePreset = CommandSet{
		{Type: msgTypeOne, Payload: []byte{0x07, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Type: msgTypeOne, Payload: []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

// SelectPreset will change to one of the 5 user presets.
func SelectPreset(channelNumber int8) CommandSet {
	if channelNumber > 5 {
		channelNumber = 0
	}
	return CommandSet{
		&SingleCommand{Type: msgTypeTwo, Payload: []byte{0x0b, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, singleByteInt(channelNumber), 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
}
