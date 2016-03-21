package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	host string
	conn net.Conn
}

type result struct {
	response string
	err      error
}

func NewClient(host string) (*Client, error) {
	client := &Client{
		host: host,
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", host)

	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return nil, err
	}

	client.conn = conn

	return client, nil
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) Get(key string) (string, error) {
	command := "GET " + key
	r := client.sendRequest(command)

	if r.err != nil {
		return "", r.err
	}

	switch r.response {
	case "emptycommand", "wrongcommand", "protoerr":
		return "", nil
	default:
		return string(r.response), nil
	}

	panic("unreachable")
}

func (client *Client) Set(key, value string, ttl ...int) (bool, error) {
	command := "SET " + key + " " + value

	if len(ttl) > 0 {
		command += " " + strconv.Itoa(ttl[0])
	}

	r := client.sendRequest(command)

	if r.err != nil {
		return false, r.err
	}

	switch r.response {
	case "emptycommand", "wrongcommand", "protoerr":
		return false, nil
	default:
		return string(r.response) == value, nil
	}

	panic("unreachable")
}

func (client *Client) Incr(key string, delta ...int) (string, error) {
	command := "INCR " + key

	if len(delta) > 0 {
		command += " " + strconv.Itoa(delta[0])
	}

	r := client.sendRequest(command)

	if r.err != nil {
		return "", r.err
	}

	switch r.response {
	case "emptycommand", "wrongcommand", "protoerr":
		return "", nil
	default:
		return string(r.response), nil
	}

	panic("unreachable")
}

func (client *Client) Decr(key string, delta ...int) (string, error) {
	command := "DECR " + key

	if len(delta) > 0 {
		command += " " + strconv.Itoa(delta[0])
	}

	r := client.sendRequest(command)

	if r.err != nil {
		return "", r.err
	}

	switch r.response {
	case "emptycommand", "wrongcommand", "protoerr":
		return "", nil
	default:
		return string(r.response), nil
	}

	panic("unreachable")
}

func (client *Client) Delete(key string) (bool, error) {
	command := "DEL " + key

	r := client.sendRequest(command)

	if r.err != nil {
		return false, r.err
	}

	switch r.response {
	case "emptycommand", "wrongcommand", "protoerr":
		return false, nil
	default:
		return (string(r.response) == "true"), nil
	}

	panic("unreachable")
}

func (client *Client) Add(key, value string, ttl ...int) (bool, error) {
	command := "ADD " + key + " " + value

	if len(ttl) > 0 {
		command += " " + strconv.Itoa(ttl[0])
	}

	r := client.sendRequest(command)

	if r.err != nil {
		return false, r.err
	}

	switch r.response {
	case "emptycommand", "wrongcommand", "protoerr":
		return false, nil
	default:
		return (r.response == value), nil
	}

	panic("unreachable")
}

func (client *Client) sendRequest(command string) *result {
	r := &result{}

	_, err := fmt.Fprintf(client.conn, command+"\n")

	if err != nil {
		r.err = err
		return r
	}

	response, err := bufio.NewReader(client.conn).ReadString('\n')

	if err != nil {
		r.err = err
		return r
	}

	r.response = strings.TrimSpace(string(response))

	return r
}
