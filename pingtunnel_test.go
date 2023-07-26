package pingtunnel

import (
	"fmt"
	"testing"

	"github.com/esrrhs/pingtunnel/pkg/msg"
	"google.golang.org/protobuf/proto"
)

func Test0001(t *testing.T) {

	my := &msg.PingMsg{}
	my.Id = "12345"
	my.Target = "111:11"
	my.Type = 12
	my.Data = make([]byte, 0)
	dst, _ := proto.Marshal(my)
	fmt.Println("dst = ", dst)

	my1 := &msg.PingMsg{}
	proto.Unmarshal(dst, my1)
	fmt.Println("my1 = ", my1)
	fmt.Println("my1.Data = ", my1.Data)

	proto.Unmarshal(dst[0:4], my1)
	fmt.Println("my1 = ", my1)

}
