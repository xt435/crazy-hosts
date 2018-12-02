package sockcomm

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

func connection(ip string, port int) {
	for {
		time.Sleep(time.Millisecond * 1000 * 3)
		con, err := conn(ip, port)
		if err != nil {
			fmt.Println(err)
		} else {
			finalConn = con
			break
		}
	}
	_, errWrite := finalConn.Write([]byte("8ff7d7af11b2c16c7a2f42e12fe232e233840442" + time.Now().Format(time.RFC850)))
	if errWrite != nil {
		fmt.Println(errWrite)
		return
	}
	connBuf := bufio.NewReader(finalConn)
	go recvSockCli(connBuf)
	for {
		_, errWrite := finalConn.Write([]byte("8ff7d7af11b2c16c7a2f42e12fe232e233840442" + time.Now().Format(time.RFC850)))
		if errWrite != nil {
			fmt.Println(errWrite)
			time.Sleep(time.Millisecond * 1000 * 3)
			continue
		}
		time.Sleep(time.Millisecond * 1000 * 3)
	}
}

func recvSockCli(connBuf *bufio.Reader) {
	for {
		mess, err := connBuf.ReadString('\n')
		if err != nil {
			fmt.Println("conn err==")
			fmt.Println(err)
			time.Sleep(time.Millisecond * 1000 * 1)
			return
		}
		time.Sleep(time.Millisecond * 1000 * 1)
		if len(mess) > 0 {
			fmt.Println("mess==" + mess)
		}
	}
}

func conn(ip string, port int) (net.Conn, error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	return conn, err
}

func MainClient(ip string, port int) {
	fmt.Println("time init==" + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
	checkLocal, checkLocalPort := GetOutboundIP()
	fmt.Println("LOCAL::" + checkLocal.String() + ":" + strconv.Itoa(checkLocalPort))
	connection(ip, port)
}

func GetOutboundIP() (net.IP, int) {
	conn, err := net.Dial(NETWORK_UDP, LOCAL_OUTBOUND_CHECK_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, localAddr.Port
}
