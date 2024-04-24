package outbound

// Packet Structure is as follows:
// 0: Command
// 1-2: length (i.e. 46)
// 3: Sys Info flag, always 0xCC (Do not send video card info)
// 4-5: Number of servers
// 6-7: Server Index (Always 0)
// 8-39: Server Name
// 40: Percentage full
// 41: Timezone
// 42-45: IP Address (Backwards)

import "net"

const (
	PacketSize             = 46
	CmdGameServerList byte = 0xA8
)

type GameServerListMessage struct {
	data []byte
}

func NewGameServerListMessage(name string, ip net.IPAddr) *GameServerListMessage {
	d := make([]byte, PacketSize)
	d[0] = CmdGameServerList
	d[1] = 0x00
	d[2] = 0x2E
	d[3] = 0xCC
	d[4] = 0x00
	d[5] = 0x01
	d[6] = 0x00
	copy(d[8:39], name)
	d[40] = 0x00
	d[41] = 0x00
	copy(d[42:46], ip.IP)
	return &GameServerListMessage{
		data: d,
	}
}
