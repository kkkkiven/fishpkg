// Copyright (c) 2020. Homeland Interactive Technology Ltd. All rights reserved.

package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/kkkkiven/fishpkg/util"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

// ConnTimeoutSeconds 连接超时秒数
var ConnTimeoutSeconds = 5

// SSHConfig SSH 代理配置
type SSHConfig struct {
	Host string `yaml:"host" json:"host"` // SSH 地址
	Port string `yaml:"port" json:"port"` // SSH 端口
	User string `yaml:"user" json:"user"` // SSH 用户名
	Pass string `yaml:"pass" json:"pass"` // SSH 密码
}

// Config 配置参数
type Config struct {
	Host            string     `yaml:"host" json:"host"`                           // 数据库地址
	Port            string     `yaml:"port" json:"port"`                           // 数据库端口
	Dbname          string     `yaml:"dbname" json:"dbname"`                       // 数据库库名
	User            string     `yaml:"user" json:"user"`                           // 数据库用户名
	Pass            string     `yaml:"pass" json:"pass"`                           // 数据库密码
	Charset         string     `yaml:"charset" json:"charset"`                     // 数据库字符集
	MaxIdle         int        `yaml:"max_idle" json:"max_idle"`                   // 最大闲置连接数
	MaxConn         int        `yaml:"max_conn" json:"max_conn"`                   // 数据库最大连接数
	ConnMaxLifetime int        `yaml:"conn_max_lifetime" json:"conn_max_lifetime"` // 数据库连接生命周期, 单位: 秒
	SSHAgent        *SSHConfig `yaml:"sshagent" json:"sshagent"`                   // SSH 代理配置
}

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

// ConnMySQL 连接MySQL数据库, 并将其设置为包对默认连接对象
func ConnMySQL(cfg *Config) error {
	i, err := GetMySQLConnInstance(cfg)
	if err != nil {
		return err
	}
	DefaultDB = i
	return nil
}

// GetMySQLConnInstance 获取一个MySQL连接实例
func GetMySQLConnInstance(cfg *Config) (*Database, error) {
	var err error
	ins, network := NewDatabase(), "tcp"

	if cfg.SSHAgent != nil {
		sshClient, err := DialSSHWithPassword(
			cfg.SSHAgent.Host+":"+cfg.SSHAgent.Port,
			cfg.SSHAgent.User, cfg.SSHAgent.Pass)

		if err != nil {
			return nil, err
		}
		mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{sshClient}).Dial)
		network = "mysql+tcp"
		ins.SSHClient = sshClient
	}

	ins.DB, err = sql.Open("mysql",
		cfg.User+":"+cfg.Pass+
			"@"+network+"("+cfg.Host+":"+cfg.Port+")/"+
			cfg.Dbname+"?charset="+cfg.Charset+"&timeout="+util.Itoa(ConnTimeoutSeconds)+"s")

	if err != nil {
		return nil, fmt.Errorf("%s (%s@%s:%s/%s)", err.Error(), cfg.User, cfg.Host, cfg.Port, cfg.Dbname)
	}

	// 设置Mysql闲置和最大连接数
	// 将闲置连接设置为0是因为连接如果长时间保持,
	// 可能会因为Mysql本身设置有wait_timeout而单方面断开,
	// 此时将发生错误 packets.go:32: unexpected EOF...
	// 所以此处将最大闲置数设置0让其使用后自动断开
	if cfg.MaxIdle > 0 {
		ins.DB.SetMaxIdleConns(cfg.MaxIdle)
		ins.DB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}
	if cfg.MaxConn > 0 {
		ins.DB.SetMaxOpenConns(cfg.MaxConn)
	}
	err = ins.DB.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s (%s@%s:%s/%s)", err.Error(), cfg.User, cfg.Host, cfg.Port, cfg.Dbname)
	}
	return ins, nil
}

// DialSSHWithPassword 通过密码拨号 SSH 客户端
func DialSSHWithPassword(addr, user, pass string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if pass != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(pass),
		}
	}

	return DialSSH("tcp", addr, config)
}

// DialSSHWithKey 通过密钥拨号 SSH 客户端
func DialSSHWithKey(addr, user, keyFile string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return DialSSH("tcp", addr, config)
}

// DialSSH 拨号 SSH 客户端
func DialSSH(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	client, err := ssh.Dial(network, addr, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
