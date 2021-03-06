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

// 获取dns记录
func GetDnsRecord(domain string) string {
	domainSubs := strings.Split(domain, ".")
	for i, j := 0, len(domainSubs)-1; i < j; i, j = i+1, j-1 {
		domainSubs[i], domainSubs[j] = domainSubs[j], domainSubs[i]
	}

	key := fmt.Sprintf("%s/%v", prefix, strings.Join(domainSubs, "/"))
	res, err := Get(key)
	if err != nil || len(res) <= 0 {
		domainSubs[len(domainSubs)-1] = "*"
		key = fmt.Sprintf("%s/%v", prefix, strings.Join(domainSubs, "/"))
		res, err := Get(key)
		if err != nil || len(res) <= 0 {
			return ""
		}
		return res["host"]
	}
	return res["host"]
}

// 获取所有dns记录
func GetAllDnsRecord() map[string]string {
	result := map[string]string{}
	res, err := GetWithPrefix(prefix)
	if err != nil {
		EtcdLogger.Error(err.Error())
		return result
	}
	for k, v := range res {
		v = strings.TrimSpace(v)
		hostData := CoreDnsHost{}
		json.Unmarshal([]byte(v), &hostData)
		domainSubs := strings.Split(strings.Replace(k, fmt.Sprintf("%s/", prefix), "", -1), "/")
		for i, j := 0, len(domainSubs)-1; i < j; i, j = i+1, j-1 {
			domainSubs[i], domainSubs[j] = domainSubs[j], domainSubs[i]
		}
		result[strings.Join(domainSubs, ".")] = hostData.Host
	}
	return result
}
