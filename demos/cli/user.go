package main

import (
	"fmt"
	"sync"

	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
)

type User struct {
	*client.CLBUser
}

type CustomClient struct {
	sync.RWMutex
	client.Client

	user *User
	ui   *CLUI
}

func NewUser(opts client.CLBOptions) *User {
	ncl, ok := client.NewCLBUser(opts)
	if !ok {
		return nil
	}
	ret := &User{ncl}
	return ret
}

func NewCustomClient(uid, pid, cid string) *CustomClient {
	return &CustomClient{sync.RWMutex{}, client.Client{Username: uid, Password: pid, ClientId: cid}, nil, nil}
}

func (self *CustomClient) Connected(opts protobase.OptionInterface) bool {
	logger.Infof("+ [USER] %s connected with opts %+v.\n", self.Username, opts.(*protocol.ConnackOpts))
	self.ui.print(fmt.Sprintf("+ Connected to broker.\n"), "input")
	self.user.SetConnected(true)
	self.user.Conn.Subscribe(*sr, byte(*qos), pcallback)
	return true
}

func (self *CustomClient) Disconnected(opts protobase.OptCode) {
	logger.Infof("+ [USER] %s disconnected.\n", self.Username)
	self.ui.print(fmt.Sprintf("- Disconnected from broker.\n"), "input")
	self.user.SetConnected(true)
}

func (self *CustomClient) Subscribe(msg protobase.MsgInterface) {
	logger.Infof("+ [USER] %s subscribed to %s.\n", self.Username, msg.Envelope().Route())
	self.ui.print(fmt.Sprintf("+ Subscribed to %s .\n", msg.Envelope().Route()), "input")
}

func (self *CustomClient) Publish(msg protobase.MsgInterface) {
	var (
		dir      protobase.MsgDir               = msg.Dir()
		envelope protobase.MsgEnvelopeInterface = msg.Envelope()
	)

	var (
		topic   string = envelope.Route()
		message []byte = envelope.Payload()
	)
	switch dir {
	case protobase.MDInbound:
		logger.Infof("+ [USER][publish] %s is sending to topic [%s], message [%s].\n", self.Username, topic, string(message))
	case protobase.MDOutbound:
		logger.Infof("+ [USER][publish] %s has receive dtopic [%s], message [%s].\n", self.Username, topic, string(message))
	default:
		logger.Infof("- [USER][publish] unknown direction flag(%d).", int(dir))
	}
	self.ui.print(fmt.Sprintf("> %s\n", string(message)), "input")
}

func (self *CustomClient) pubCli(msg []byte) {
	// self.ui.print("pubcli\n", "input")
	self.user.Conn.Publish(*sr, msg, byte(*qos), self.pubCallback)
}

func (self *CustomClient) pubCallback(opts protobase.OptionInterface, msg protobase.MsgInterface) {
	// TODO
}
