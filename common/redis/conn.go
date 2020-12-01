//
//	包含redis连接的方法
//
// 	Redis的连接方式
// 	1. 普通连接
// 	2. SSH普通连接
// 	3. 连接池连接
// 	4. SSH连接池连接
//
package redis

import (
	"fmt"
	"net"
	_ "strings"
	_ "time"

	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/ssh"
)

// getSSHClient 连接ssh
// addr : 主机地址, 如: 127.0.0.1:22
// user : 用户
// pass : 密码
// 返回 ssh连接
func getSSHClient(user string, pass string, addr string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshConn, err := net.Dial("tcp", addr)
	if nil != err {
		fmt.Println("net dial err: ", err)
		return nil, err
	}

	clientConn, chans, reqs, err := ssh.NewClientConn(sshConn, addr, config)
	if nil != err {
		sshConn.Close()
		fmt.Println("ssh client conn err: ", err)
		return nil, err
	}

	client := ssh.NewClient(clientConn, chans, reqs)
	fmt.Println("ssh client = ", user, pass, addr, client)
	return client, nil
}

// RConn 普通连接
// ip : redis服务地址
// port:  Redis 服务端口
// password  Redis 服务密码
// 返回redis连接
func RConn(ip string, port int, password string) (redis.Conn, error) {
	host := fmt.Sprintf("%s:%d", ip, port)
	conn, err := redis.Dial("tcp", host)
	if nil != err {
		fmt.Println("dial to redis addr err: ", err)
		return nil, err
	}

	fmt.Println(conn)

	//TODO: 针对无密码的连接
	if password == "" {
		return conn, nil
	}

	if _, authErr := conn.Do("AUTH", password); authErr != nil {
		fmt.Println("redis auth password error: ", authErr)
		return nil, fmt.Errorf("redis auth password error: %s", authErr)
	}

	return conn, nil
}

// RSSHConn SSH普通连接
// addr : SSH主机地址, 如: 127.0.0.1:22
// user : SSH用户
// pass : SSH密码
// ip : redis服务地址
// port:  Redis 服务端口
// password  Redis 服务密码
// 返回redis连接
func RSSHConn(user string, pass string, addr string, ip string, port int, password string) (redis.Conn, error) {
	sshClient, err := getSSHClient(user, pass, addr)
	if nil != err {
		fmt.Println(err)
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", ip, port)
	conn, err := sshClient.Dial("tcp", host)
	if nil != err {
		fmt.Println("dial to redis addr err: ", err)
		return nil, err
	}

	redisConn := redis.NewConn(conn, -1, -1)

	//TODO: 针对无密码的连接
	if password == "" {
		return redisConn, nil
	}

	if _, authErr := redisConn.Do("AUTH", password); authErr != nil {
		fmt.Println("redis auth password error: ", authErr)
		return nil, fmt.Errorf("redis auth password error: %s", authErr)
	}

	return redisConn, nil
}

//指定db的连接
func SelectDB(rc redis.Conn, dbnumber int64) (redis.Conn, error) {
	_, err := rc.Do("select", fmt.Sprintf("%d", dbnumber))
	if err != nil {
		fmt.Println("redis select db error: ", err)
	}
	return rc, err
}
