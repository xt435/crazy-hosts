package monitorspack

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	//NETWORKUDP proto udp
	NETWORKUDP string = "udp"
	//NETWORKTCP proto tcp
	NETWORKTCP string = "tcp"
	//LocalOutboundCheckAddress for outbound socket check
	LocalOutboundCheckAddress = "8.8.8.8:80"
)

//ChancesForRetryMax max chances for connecting. default 10
var ChancesForRetryMax = 10
var bytesRepo []map[string][]string
var finalConn net.Conn
var wg sync.WaitGroup

func setMaxRetry(numberMax int) {
	ChancesForRetryMax = numberMax
}

//Conn the connection parameters
type Conn struct {
	IP   string
	Port int
	Sq   chan string
	Rq   chan string
}

func connection(conner Conn) {
	var chancesForRetry = 0
	for {
		con, err := conn(conner.IP, conner.Port)
		if err != nil {
			fmt.Println(err)
			chancesForRetry++
			if chancesForRetry == ChancesForRetryMax-1 {
				fmt.Printf("Cannot connect to %s : %d after %d tries", conner.IP, conner.Port, ChancesForRetryMax)
				break
			}
			time.Sleep(time.Millisecond * 3000) // retry after 3 secs
		} else {
			finalConn = con
			break
		}
	}
	_, errWriteFirst := finalConn.Write([]byte("are we in?\n" + time.Now().Format(time.RFC3339)))
	if errWriteFirst != nil {
		fmt.Println(errWriteFirst)
		fmt.Println("About To Reconn...")
		connection(conner)
		return
	}
	connBuf := bufio.NewReader(finalConn)
	go recvSockCli(connBuf, conner)
	for {
		m := <-conner.Sq
		if len(m) > 0 {
			fmt.Printf("mout==%s\n", m)
			_, errWrite := finalConn.Write([]byte(m))
			if errWrite != nil {
				fmt.Println(errWrite)
				time.Sleep(time.Millisecond * 300)
				fmt.Println("About To Reconn...")
				connection(conner)
				return
			}
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func recvSockCli(connBuf *bufio.Reader, conner Conn) {
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
			conner.Rq <- mess
		}
	}
}

func conn(ip string, port int) (net.Conn, error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	return conn, err
}

// MainClient sender and receiver
func MainClient(mon []MonContent) {
	fmt.Printf("time-init==%s = %s\n", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10), time.RFC3339)
	checkLocal, checkLocalPort := GetOutboundIP()
	fmt.Printf("LOCAL::%s\n", checkLocal.String()+":"+strconv.Itoa(checkLocalPort))
	if len(mon) > 0 {
		for i := range mon {
			conn := Conn{IP: mon[i].IP, Port: mon[i].Port, Sq: mon[i].Send, Rq: mon[i].Recv}
			// for {
			// 	mess := <-sq
			// 	if len(mess) > 0 {
			// 		fmt.Printf("mosters::%s\n", mess)
			// 		rq <- "verywell"
			// 	}
			// 	time.Sleep(time.Millisecond * 300)
			// }
			connection(conn)
		}
	}

}

// GetOutboundIP check self
func GetOutboundIP() (net.IP, int) {
	conn, err := net.Dial(NETWORKUDP, LocalOutboundCheckAddress)
	if err != nil {
		log.Fatal(err)
		fmt.Println("!! WARNING !! no connection to the internet")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, localAddr.Port
}
