package main

import (
	"flag"
	"fmt"
	"net"
	"regexp"
	"strings"
)

//go env -w CGO_ENABLED=0 -w GOOS=android -w GOARCH=arm  build -ldflags "-s -w" get_public_ipv6.go

func main() {
	public := flag.Bool("public", false, "显示“临时的”公网 IPv6，"+
		"该 IPv6 地址是路由器临时分配的，会在一定时间刷新，也是你上网时对外暴露的 IPv6 地址。"+
		"此地址重新连接网络后，就会改变")
	permanent := flag.Bool("permanent", false, "显示“长久的”公网 IPv6，"+
		"该 IPv6 地址后面几段永远是相同的，这是根据你网卡 MAC 地址生成的，不建议暴露此 IPv6 地址。"+
		"此地址即使重新连接网络后，也不会改变")

	var result string

	flag.Parse()

	if *public == true {
		result = public_ipv6()
	} else if *permanent == true {
		result = permanent_ipv6()
	} else {
		result = "public_ipv6:>>" + public_ipv6()
		result += "\npermanent_ipv6:>>" + permanent_ipv6()
	}

	fmt.Println(result)
}

/**
 * 获取“永久”公网 IPv6
 *
 * @return string 如果获取成功将返回类似 2400:3200::1 的 IPv6 地址
 *				  如果不存在 IPv6 地址或获取失败，则返回值为空字符
 */
func public_ipv6() string {
	//[2400:3200::1] 是 “阿里 IPv6 DNS”
	conn, _ := net.Dial("udp", "[2400:3200::1]:5353")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	ipv6 := localAddr[0:idx]

	if ipv6 != "" {
		//假设获取到的结果是 [2400:3200::1]
		//删除获取到 IPv6 地址的修饰符 []
		//使它变成 2400:3200::1
		ipv6 = strings.Replace(ipv6, "[", "", 1)
		ipv6 = strings.Replace(ipv6, "]", "", 1)
	}

	return ipv6
}

/**
 * 获取“永久”公网 IPv6
 *
 * @return string 如果获取成功将返回类似 2400:3200::1 的 IPv6 地址
 *				  如果不存在 IPv6 地址或获取失败，则返回值为空字符
 */
func permanent_ipv6() string {
	s, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range s {
		i := regexp.MustCompile(`(\w+:){7}\w+`).FindString(a.String())
		if strings.Count(i, ":") == 7 {
			return i
		}
	}
	return ""
}

//SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 go build -ldflags "-s -w" get_public_ipv6.go
