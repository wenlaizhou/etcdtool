package etcdtool

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/wenlaizhou/middleware"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"strings"
	"time"
)

var endpoints []string
var timeout time.Duration
var tlsConf *tls.Config

var EtcdLogger = middleware.GetLogger("etcd")

// 获取etcd终端
func GetEndpoints() string {
	return strings.Join(endpoints, ",")
}

// 初始化配置
func InitConf(eps []string, timeOut int, keyPath string, caPath string, rootCaPath string) error {
	if len(eps) <= 0 {
		return errors.New("etcd配置错误, endpoints地址为空")
	}
	for _, endpoint := range eps {
		endpoint = strings.TrimSpace(endpoint)
		if len(endpoint) > 0 {
			endpoints = append(endpoints, endpoint)
		}
	}
	timeout = time.Duration(timeOut) * time.Second

	if middleware.HasEmptyString(keyPath, caPath, rootCaPath) {
		tlsConf = nil
		return nil
	}

	cert, err := tls.LoadX509KeyPair(caPath, keyPath)
	if err != nil {
		return err
	}

	caData, err := ioutil.ReadFile(rootCaPath)
	if err != nil {
		return err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	tlsConf = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return nil

}

// 获取客户端
func GetClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
		TLS:         tlsConf,
	})
}

// get
func Get(key string) (map[string]string, error) {
	cli, err := GetClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := cli.Get(ctx, key)
	cancel() // todo explain cancel
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Count > 0 {
		res := make(map[string]string)
		for _, kv := range resp.Kvs {
			res[string(kv.Key)] = string(kv.Value)
		}
		return res, nil
	}
	return map[string]string{}, nil
}

// get with prefix
func GetWithPrefix(key string) (map[string]string, error) {
	cli, err := GetClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	cancel() // todo explain cancel
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Count > 0 {
		res := make(map[string]string)
		for _, kv := range resp.Kvs {
			res[string(kv.Key)] = string(kv.Value)
		}
		return res, nil
	}
	return map[string]string{}, nil
}

// put
func Put(key string, value string) error {
	cli, err := GetClient()
	if err != nil {
		return err
	}
	defer cli.Close()
	_, err = cli.Put(context.TODO(), key, value)
	return err
}

// delete
func Delete(key string) error {
	cli, err := GetClient()
	if err != nil {
		return err
	}
	defer cli.Close()
	_, err = cli.Delete(context.TODO(), key)
	return err
}

// deleteAll
func DeleteAll(key string) error {
	cli, err := GetClient()
	if err != nil {
		return err
	}
	defer cli.Close()
	_, err = cli.Delete(context.TODO(), key, clientv3.WithPrefix())
	return err
}

// 获取所有key
func GetAllKeys() ([]string, error) {
	cli, err := GetClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := cli.Get(ctx, "", clientv3.WithKeysOnly(), clientv3.WithPrefix())
	cancel() // todo explain cancel
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Count > 0 {
		var res []string
		for _, kv := range resp.Kvs {
			res = append(res, string(kv.Key))
		}
		return res, nil
	}
	return nil, nil

}
