package etcdtool

import (
	"encoding/json"
	"fmt"
	"strings"
)

const prefix = "/skydns"

const defaultTtl = 60

type CoreDnsHost struct {
	Host string `json:"host"`
	Ttl  int    `json:"ttl"`
}

// 添加一条dns记录
//
// dig +short www.x.com
func AddDnsRecord(domain string, ip string) error {
	domainSubs := strings.Split(domain, ".")
	for i, j := 0, len(domainSubs)-1; i < j; i, j = i+1, j-1 {
		domainSubs[i], domainSubs[j] = domainSubs[j], domainSubs[i]
	}
	host := CoreDnsHost{
		Host: ip,
		Ttl:  defaultTtl,
	}
	jsonData, _ := json.Marshal(host)
	EtcdLogger.InfoF("新增一条dns记录: %s, %s", domain, ip)
	key := fmt.Sprintf("%s/%v", prefix, strings.Join(domainSubs, "/"))
	return Put(key, string(jsonData))
}
