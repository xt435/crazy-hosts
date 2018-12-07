package monitorspack

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	NETWORK_UDP = "udp"
	NETWORK_TCP = "tcp"

	LOCAL_OUTBOUND_CHECK_ADDRESS = "8.8.8.8:80"
)

var bytesRepo []map[string][]string
var finalConn net.Conn

func connection(ip string, port int, sq []string, rq []string) {
	for {
		con, err := conn(ip, port)
		if err != nil {
			fmt.Println(err)
		} else {
			finalConn = con
			break
		}
	}
	_, errWriteFirst := finalConn.Write([]byte("are we in?\n" + time.Now().Format(time.RFC3339)))
	if errWriteFirst != nil {
		fmt.Println(errWriteFirst)
		return
	}
	connBuf := bufio.NewReader(finalConn)
	go recvSockCli(connBuf, rq)
	for {
		if len(sq) > 0 {
			head := sq[0]
			sq = sq[1:]
			_, errWrite := finalConn.Write([]byte(head))
			if errWrite != nil {
				fmt.Println(errWrite)
				time.Sleep(time.Millisecond * 300)
				continue
			}
		}
		time.Sleep(time.Millisecond * 300)
	}
}

func recvSockCli(connBuf *bufio.Reader, rq []string) {
	for {
		mess, err := connBuf.ReadString('\n')
		if err != nil {
			fmt.Println("conn err==")
			fmt.Println(err)
			time.Sleep(time.Millisecond * 300)
			return
		}
		time.Sleep(time.Millisecond * 30)
		if len(mess) > 0 {
			fmt.Printf("mess_raw==%s\n", mess)
			rq = append(rq, mess)
		}
	}
}

func conn(ip string, port int) (net.Conn, error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	return conn, err
}

// MainClient sender and receiver
func MainClient(ip string, port int, sq []string, rq []string) {
	fmt.Printf("time-init==%s", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
	checkLocal, checkLocalPort := GetOutboundIP()
	fmt.Printf("LOCAL::%s", checkLocal.String()+":"+strconv.Itoa(checkLocalPort))
	connection(ip, port, sq, rq)
}

// GetOutboundIP check self
func GetOutboundIP() (net.IP, int) {
	conn, err := net.Dial(NETWORK_UDP, LOCAL_OUTBOUND_CHECK_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, localAddr.Port
}
