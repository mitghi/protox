package protocol

import (
	"github.com/mitghi/protox/protobase"
)

// TODO:
// . embed utils into []byte as protodata
// . refactor payload fields into separate struct

// QueueAck constants
const (
	QAcNone byte = iota
	QAcOK
	QAcERR
)

// Queue constants
const (
	QAInitialize protobase.QAction = iota
	QADestroy
	QADrain
	QANone
)

var (
	_ protobase.EDProtocol = (*Connect)(nil)
)

// Connect establish connection.
type Connect struct {
	Protocol

	ClientId   string
	Username   string
	Password   string
	Version    string
	KeepAlive  int
	CleanStart bool
}

// Connack is a control packet. It acknowledges the incomming connections
// and includes a `ResultCode` which determines connection status and an
// optional `SessionId` which should be used by the client for resuming
// previous states.
type Connack struct {
	Protocol

	ResultCode byte
	SessionId  string
	// TODO
	// . add config fields from broker to client
}

// Disconnect is a control packet. It temrinates the connection.
type Disconnect struct {
	Protocol
}

type Ping struct {
	Protocol
}

type Pong struct {
	Protocol
}

//
type Puback struct {
	Protocol
}

//
type Suback struct {
	Protocol
}

//
type Unsuback struct {
	Protocol
}

//
type Publish struct {
	Protocol

	Topic   string
	Message []byte
}

type QAck struct {
	Protocol

	Code byte
}

type Queue struct {
	Protocol

	Message    []byte
	Mark       []byte
	Address    string
	ReturnPath string
	Action     protobase.QAction
}

//
type Subscribe struct {
	Protocol

	Topic string
}

//
type UnSubscribe struct {
	Protocol

	Topic string
}

// NewConnect returns a new `Connect`
// control packet.
func NewConnect(packet protobase.PacketInterface) (c *Connect) {
	c = &Connect{
		Protocol: NewProtocol(protobase.CCONNECT),
		// TODO
		// . add control byte options ( 0x0f )
	}
	if err := c.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return c
}

// NewDisconnect returns a new `Disconnect` control packet.
func NewDisconnect(packet protobase.PacketInterface) (d *Disconnect) {
	d = &Disconnect{
		Protocol: NewProtocol(protobase.CDISCONNECT),
	}
	if err := d.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return d
}

// NewConnack returns a new `Connack` control packet.
func NewConnack(packet protobase.PacketInterface) (ca *Connack) {
	ca = &Connack{
		Protocol: NewProtocol(protobase.CCONNACK),
	}
	if err := ca.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return ca
}

// NewPing returns a new Ping control packet. It is not the responsibility of broker to send
// ping control packets ( it is for client ).
func NewPing(packet protobase.PacketInterface) (p *Ping) {
	p = &Ping{
		Protocol: NewProtocol(protobase.CPING),
	}
	if err := p.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return p
}

// NewPong returns a pointer to a new `Pong` packet.
func NewPong(packet protobase.PacketInterface) (p *Pong) {
	p = &Pong{
		Protocol: NewProtocol(protobase.CPONG),
	}
	if err := p.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return p
}

//
func NewPuback(packet protobase.PacketInterface) (pa *Puback) {
	pa = &Puback{
		Protocol: NewProtocol(protobase.CPUBACK),
	}
	if err := pa.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return pa
}

//
func NewPublish(packet protobase.PacketInterface) (p *Publish) {
	p = &Publish{
		Protocol: NewProtocol(protobase.CPUBLISH),
		Topic:    "",
	}
	if err := p.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return p
}

func NewQAck(packet protobase.PacketInterface) (qa *QAck) {
	qa = &QAck{
		Protocol: NewProtocol(protobase.CQUEUEACK),
		Code:     QAcNone,
	}
	if err := qa.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return qa
}

// - MARK: Initializers.

func NewQueue() *Queue {
	return &Queue{
		Protocol: NewProtocol(protobase.CQUEUE),
		Action:   QANone,
	}
}

//
func NewSuback(packet protobase.PacketInterface) (sa *Suback) {
	sa = &Suback{
		Protocol: NewProtocol(protobase.CSUBACK),
	}
	if err := sa.DecodeFrom(packet.GetData()); err != nil {
		return nil
	}
	return sa
}

//
func NewSubscribe(packet protobase.PacketInterface) (s *Subscribe) {
	s = &Subscribe{
		Protocol: NewProtocol(protobase.CSUBSCRIBE),
	}
	if err := s.DecodeFrom(packet.GetData()); err != nil {
		logger.FDebug("onPUBACK", "- [PubAck] uanble to decode .", "error", err)
		return
	}
	return s
}

//
func NewUnsuback(packet protobase.PacketInterface) (usa *Unsuback) {
	usa = &Unsuback{
		Protocol: NewProtocol(protobase.CUNSUBACK),
	}
	if err := usa.DecodeFrom(packet.GetData()); err != nil {
		logger.FDebug("onPUBACK", "- [PubAck] uanble to decode .", "error", err)
		return
	}
	return usa
}

//
func NewUnSubscribe(packet protobase.PacketInterface) (us *UnSubscribe) {
	us = &UnSubscribe{
		Protocol: NewProtocol(protobase.CUNSUBSCRIBE),
	}
	if err := us.DecodeFrom(packet.GetData()); err != nil {
		logger.FDebug("onPUBACK", "- [PubAck] uanble to decode .", "error", err)
		return
	}
	return us
}

// RAW

// NewConnect returns a new `Connect`
// control packet.
func NewRawConnect() *Connect {
	return &Connect{
		Protocol: NewProtocol(protobase.CCONNECT),
		// TODO
		// . add control byte options ( 0x0f )
	}
}

// NewDisconnect returns a new `Disconnect` control packet.
func NewRawDisconnect() *Disconnect {
	result := &Disconnect{
		Protocol: NewProtocol(protobase.CDISCONNECT),
	}

	return result
}

// NewConnack returns a new `Connack` control packet.
func NewRawConnack() *Connack {
	result := &Connack{
		Protocol:   NewProtocol(protobase.CCONNACK),
		SessionId:  "",
		ResultCode: 0x0,
	}

	return result
}

// NewPing returns a new Ping control packet. It is not the responsibility of broker to send
// ping control packets ( it is for client ).
func NewRawPing() *Ping {
	result := &Ping{
		Protocol: NewProtocol(protobase.CPING),
	}
	return result
}

// NewPong returns a pointer to a new `Pong` packet.
func NewRawPong() *Pong {
	result := &Pong{
		Protocol: NewProtocol(protobase.CPONG),
	}

	return result
}

//
func NewRawPuback() *Puback {
	return &Puback{
		Protocol: NewProtocol(protobase.CPUBACK),
	}
}

//
func NewRawPublish() *Publish {
	return &Publish{
		Protocol: NewProtocol(protobase.CPUBLISH),
		Topic:    "",
	}
}

func NewRawQAck() *QAck {
	return &QAck{
		Protocol: NewProtocol(protobase.CQUEUEACK),
		Code:     QAcNone,
	}
}

// - MARK: Initializers.

func NewRawQueue() *Queue {
	return &Queue{
		Protocol: NewProtocol(protobase.CQUEUE),
		Action:   QANone,
	}
}

//
func NewRawSuback() *Suback {
	return &Suback{
		Protocol: NewProtocol(protobase.CSUBACK),
	}
}

//
func NewRawSubscribe() *Subscribe {
	return &Subscribe{
		Protocol: NewProtocol(protobase.CSUBSCRIBE),
	}
}

//
func NewRawUnsuback() *Unsuback {
	return &Unsuback{
		Protocol: NewProtocol(protobase.CUNSUBACK),
	}
}

//
func NewRawUnSubscribe() *UnSubscribe {
	return &UnSubscribe{
		Protocol: NewProtocol(protobase.CUNSUBSCRIBE),
	}
}
