package v1_0

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Client struct {
	Conn net.Conn
}

func (c *Client) Rpc(in_ string, params string) (out string, result string) {
	m := map[string]string{"in": in_, "params": params}
	request, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Fail to marshal, %s\n", err)
		return
	}
	length_prefix := make([]byte, 4)
	//length_prefix = []byte(strconv.Itoa(len(request)))
	//buf := new(bytes.Buffer)
	//binary.Write(buf, binary.LittleEndian, len(request))
	//length_prefix = buf.Bytes()
	binary.LittleEndian.PutUint32(length_prefix, uint32(len(request)))
	c.Conn.Write(length_prefix)
	c.Conn.Write(request)
	c.Conn.Read(length_prefix)
	length := binary.LittleEndian.Uint32(length_prefix)
	body := make([]byte, length)
	c.Conn.Read(body)
	response := map[string]string{}
	err2 := json.Unmarshal(body, &response)
	if err2 != nil {
		fmt.Printf("Fail to unmarshal, %s\n", err2)
		return
	}
	return response["out"], response["result"]
}
