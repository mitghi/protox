/* MIT License
*
* Copyright (c) 2018 Mike Taghavi <mitghi[at]gmail.com>
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
*/

package protobase

import (
	"net"

	"github.com/google/uuid"
)

type (
	// MsgDir is message direction. It indicates wether a message
	// is inbound or outbound.
	MsgDir byte
	// OptCode is the option type used by `OptionsInterface`.
	OptCode byte
	// ACLMode is the type for Access Control List modes.
	ACLMode byte
	// AuthMode is the type for Authenication subsystem.
	AuthMode byte
	// ClientMode is the type for various client modes such as User, Router, or Agent.
	ClientMode byte
	// AuthUserType is the type to identify user groups in a Authentication implementor.
	AuthUserType string
	// QAction is the type for identifying Queue commands.
	QAction uint
)

// CredentialsInterface is the interface for credential providers.
// It is used by authenicators.
type CredentialsInterface interface {
	GetCredentials() (string, string, string)
	IsValid() bool
	Match(CredentialsInterface) bool
	GetUID() string
	Copy() CredentialsInterface
}

// AuthInterface is the interface for authenication facilities.
type AuthInterface interface {
	HasSession(clientId string) bool
	CanAuthenticate(creds CredentialsInterface) (bool, error)
	TryAuthenticate(creds CredentialsInterface) bool
	TryUnAuthenticate(string) bool
	HasClient(string) bool
	MakeCreds(string, string, string, ...interface{}) (CredentialsInterface, error)
	Authenticate(creds CredentialsInterface) bool
	Register(creds CredentialsInterface) (result bool)
	RemoveWithIdentifier(identifier *string) (bool, error) //NOTE change this
	SetMode(AuthMode)
	GetMode() AuthMode
	GetACL() ACLInterface
	GetUserType(string) (AuthUserType, error)
}

// RetainStorageInterface is the interface for retained messages
// container.
type RetainStorageInterface interface {
	Insert([]byte, EDProtocol) error
	Find([]byte) (EDProtocol, error)
	Remove([]byte) error
}

// MsgEnvelopeInterface is the interface for content of a proto packet.
type MsgEnvelopeInterface interface {
	Route() string
	Payload() []byte
}

// MsgInterface is the interface that represents a proto packet and
// its neccessary associated meta data.
type MsgInterface interface {
	QoS() byte
	MessageId() uint16
	Dir() MsgDir
	Envelope() MsgEnvelopeInterface
	SetWishQoS(byte)
	Clone(MsgDir) MsgInterface
}

// OptionInterface is the interface that represents a option.
type OptionInterface interface {
	StateCode() OptCode
	Opts() interface{}
	Match(OptCode) bool
}

// Messages contains functionalities to store inbound and outbound packets.
// This is important for QoS > 0 levels as acknowledgments shall be sent either
// from Broker->Client or Client->Broker. It also provides UUID and persistency.

type MSGIDInterface interface {
	GetNewID(uuid.UUID) uint16
	IsOccupied(uint16) bool
	GetUUID(uint16) (uuid.UUID, bool)
	FreeId(uint16)
}

// MessageStorage is a interface that must be implemented
// in order to be passed into `ServerInterface` and
// `ProtoConnection` implementors.
type MessageStorage interface {
	AddClient(client string)
	AddInbound(client string, msg EDProtocol) bool
	AddOutbound(client string, msg EDProtocol) bool
	DeleteIn(client string, msg EDProtocol) bool
	DeleteOut(client string, msg EDProtocol) bool
	Exists(client string) bool
	GetAllOut(client string) (msgs []EDProtocol)
	GetAllOutStr(client string) (msgs []string)
	GetOutbound(string, uuid.UUID) (EDProtocol, bool)
	GetInbound(string, uuid.UUID) (EDProtocol, bool)
	GetIDStoreO(client string) MSGIDInterface
	GetIDStoreI(client string) MSGIDInterface
	Close(client string) bool
}

//
type MessageBox interface {
	AddInbound(msg EDProtocol) bool
	AddOutbound(msg EDProtocol) bool
	DeleteIn(msg EDProtocol) bool
	DeleteOut(msg EDProtocol) bool
	GetAllOut() (msgs []EDProtocol)
	GetAllOutStr() (msgs []string)
	GetOutbound(uuid.UUID) (EDProtocol, bool)
	GetInbound(uuid.UUID) (EDProtocol, bool)
	GetIDStoreO() MSGIDInterface
	GetIDStoreI() MSGIDInterface
}

// ClientInterface is the interface that must be implemented by subsystems
// which handle high level client logics. Each control packet may trigger
// a corresponding routine on structures implementing this interface. It is
// important to notice that some control packets do not trigger notifications.
type ClientInterface interface {
	Connected(OptionInterface) bool
	Disconnected(OptCode)
	Publish(MsgInterface)
	Subscribe(MsgInterface)
	GetIdentifier() string
	GetTopics() []string
	GetCreds() CredentialsInterface
	GetUser() interface{}
	SetCreds(CredentialsInterface)
	SetUser(interface{})
	Setup() error
	// TODO:
	// . investigate addition of auth mechanism
	//   e.g. SetAuthMechanism()  
}

//
type CLBUserInterface interface {
	IsRunning() bool
	IsConnected() bool
	GetExitCh() chan struct{}
	SetRunning(bool)
	SetConnected(bool)
	Setup() error
	Connect()
	Disconnect()
}

// ServerInterface is the interface that must be implemented by subsystems
// which do the serving. Each control packet will trigger a corresponding routine
// on structures implementing this interface.
type ServerInterface interface {
	NotifyDisconnected(prc ProtoConnection)
	NotifyConnected(prc ProtoConnection)
	NotifySubscribe(prc ProtoConnection, msg MsgInterface)
	NotifyPublish(prc ProtoConnection, msg MsgInterface)
	NotifyReject(prc ProtoConnection)
	NotifyQueue(prc ProtoConnection, msg MsgInterface)

	RegisterClient(prc ProtoConnection)
	Redeliver(prc ProtoConnection)
	Shutdown() (<-chan struct{}, error)

	GetStatus() uint32
	GetErrChan() <-chan struct{}
  GetStatusChan() <-chan uint32
	// TODO
	// . improve this
	// NotifyPublish(topic string, message string, prc ProtoConnection, dir MsgDir)
	// NotifySubscribe(topic string, prc ProtoConnection)
	Setup()
}

// PacketInterface is a interface to access 
// low level packet data.
type PacketInterface interface {
	SetData(*[]byte)
	SetCode(byte)
	SetLength(int)
	GetData() *[]byte
	GetCode() byte
	GetLength() int
	IsValid() bool
}

// EDProtocol is a interface used for PDU ( protocol data units ).
type EDProtocol interface {
	Encode() error
	Decode() error
	DecodeFrom(buff *[]byte) error
	String() string
	UUID() uuid.UUID
	MessageId() (bool, uint16)
	CommandCode() byte
	GetPacket() PacketInterface
	// Metadata() *ProtoMeta
	// TODO:
	// . SetCode(code byte)
	// . implement io.Reader and io.Writer
}

// ProtocolConnection is the main interface used by servers. It handles most of
// the logics neccessary for handling low level details ( such as parsing and crafting
// control packets, timeouts, send/receive and ... ).
type ProtoConnection interface {
	Handle()

	SetServer(sv ServerInterface)
	SetAuthenticator(auth AuthInterface)
	SetClientDelegate(cl func(string, string, string) ClientInterface)
	SetMessageStorage(store MessageStorage)
	SetHeartBeat(heartbeat int)
	SetInitiateTimeout(timeout int)
	SetStatus(uint32)
	SetNetConnection(net.Conn)
	SetPermissionDelegate(cl func(AuthInterface, ...string) bool)

	SendMessage(MsgInterface, bool)

	GetConnection() net.Conn
	GetClient() ClientInterface
	GetStatus() uint32
	GetErrChan() chan struct{}

	IsClean() bool
	SendRedelivery(EDProtocol)

	// TODO
	// SendPublish(MsgInterface)
}

// ProtocolClientConnection is the main interface used by clients. It handles most of
// the logics neccessary for handling low level details ( such as parsing and crafting
// control packets, timeouts, send/receive and ... ) and connects a client to the broker.
type ProtoClientConnection interface {
	Handle(PacketInterface)

	SetMessageStorage(MessageBox)
	SetHeartBeat(int)
	SetNetConnection(net.Conn)
	SetStatus(uint32)
	SetClient(ClientInterface)
	ContinueFlag(bool)

	GetConnection() net.Conn
	GetClient() ClientInterface
	GetStatus() uint32
	GetErrChan() chan struct{}
	GetTermChan() chan struct{}

	SendMessage(MsgInterface)
	SendRedelivery()

	MakeEnvelope(route string, payload []byte, qos byte, messageId uint16, dir MsgDir) MsgInterface
	Publish(string, []byte, byte, func(OptionInterface, MsgInterface)) error
	Subscribe(string, byte, func(OptionInterface, MsgInterface)) error
	Queue(QAction, string, string, []byte, []byte) error
	Disconnect() error
	// TODO
	// SetOptions(OptionInterface)
	// Publish(MsgInterface, func(OptionInterface, MsgInterface))
	// Subscribe(MsgInterface, func(OptionInterface, MsgInterface))
	// IsClean() bool
	// SendRedelivery(EDProtocol)

	// TODO:
	// . return error codes from `Publish` and `Subscribe`
}

// ClientDelegate is the signature for passing a new `ClientInterface` delegate.
type ClientDelegate func(net.Conn) ClientInterface

// ContNodeInterface is a interface for representing data nodes.
type ContNodeInterface interface {
	SetOpts(opts byte)
	SetValue(value interface{})
	GetValue()
}

// SubsProviderInterface is a interface for topic providers.
type SubsProviderInterface interface {
	searchPath([]byte, string, func(ContNodeInterface, [][]byte, int)) error
	insertPath([]byte, interface{}, string, func(ContNodeInterface)) ContNodeInterface
}

// LoggerInterface is a interface that all logging facilities
// must conform to.
type LoggerInterface interface {
	Trace(int, ...interface{})
	Tracef(int, string, ...interface{})
	Debug(string, ...interface{})
	Debugf(string, ...interface{})
	Info(string, ...interface{})
	Infof(string, ...interface{})
	Warn(string, ...interface{}) error
	Warnf(string, ...interface{}) error
	Error(string, ...interface{}) error
	Errorf(string, ...interface{}) error
	Fatal(string, ...interface{})
	Fatalf(string, ...interface{})
	Log(int, string, []interface{})
	IsDebug() bool
}

// LoggingInterface is the extended logging interface.
type LoggingInterface interface {
	LoggerInterface
	FInfo(string, string, ...interface{})
	FInfof(string, string, ...interface{})
	FDebug(string, string, ...interface{})
	FDebugf(string, string, ...interface{})
	FWarn(string, string, ...interface{}) error
	FWarnf(string, string, ...interface{}) error
	FError(string, string, ...interface{}) error
	FErrorf(string, string, ...interface{}) error
	FFatal(string, string, ...interface{})
	FFatalf(string, string, ...interface{})
	FTrace(int, string, string, ...interface{})
	FTracef(int, string, string, ...interface{})
}

// ILoggable is a interface for any object accepting an external
// logging facility.
type ILoggable interface {
	SetLogger(LoggingInterface)
}

// CLStoreInterface is a interface `ClientInterface` containers.
type CLStoreInterface interface {
	Add(string, ClientInterface)
	Get(string) ClientInterface
}

// BrokerInterface is the interface for broker implementors.
type BrokerInterface interface {
	Status() byte
}

// ACLNodeInterface is the interface used to implement access control list layered
// levels.
type ACLNodeInterface interface {
	CanDo(bool, ...string) bool
	Add(...string) error
	Unset(...string) (bool, error)
	HasIdentifier(string) bool
	HasWildIdentifier(string) bool
	MakeChild(int, string) ACLNodeInterface
	GetIdentifier(string) ACLNodeInterface
	SetValue(string, ACLNodeInterface) bool
	RemoveValue(string) bool
	Len() int
	IsResource(string) bool
}

// ACLInterface is the interface used to implement access control list.
type ACLInterface interface {
	MakeRole(string) (ACLPermInterface, error)
	GetRole(string) ACLPermInterface
	HasRole(string) bool
	GetOrCreate(string) (ACLPermInterface, bool)
}

// ACLPermInterface is the interface for individual access rules
// associated with a single entity.
type ACLPermInterface interface {
	HasPerm(string, string, string) bool
	HasExactPerm(string, string, string) bool
	SetPerm(string, string, string) error
	UnsetPerm(string, string, string) error
	SetMode(ACLMode) bool
}

// PermissionInterface is the interface that implementors must conform to
// in order to become compatible with ACL modules.
type PermissionInterface interface {
	HasPerm(string, string, string) bool
	HasExactPerm(string, string, string) bool
	SetPerm(string, string, string) error
	UnsetPerm(string, string, string) error
}

// ProtoEventInterface is the interface for
// implementing event handler responsible
// for processing protox packets. It gets
// invoked from compatible caller conforming
// to 'protobase.ConnectionState'.
type ProtoEventInterface interface {
	OnCONNECT(PacketInterface)
	OnCONNACK(PacketInterface)
	OnPUBLISH(PacketInterface)
	OnPUBACK(PacketInterface)
	OnSUBSCRIBE(PacketInterface)
	OnSUBACK(PacketInterface)
	OnPING(PacketInterface)
	OnPONG(PacketInterface)
	OnDISCONNECT(PacketInterface)
	OnQUEUE(PacketInterface)
}


// ConStateInterface is the requirement
// for implementing protox event handler.
type ConStateInterface interface {
  ProtoEventInterface
	HandleDefault(packet PacketInterface) (status bool)        // dispatch loop
	Handle(packet PacketInterface)                             // bootstrap routine
	Run()                                                      // main routine
	SetNextState()                                             // push state handler
  Shutdown()
}

// ConnectionState is the interface for status 
// of a connection. Each state must implement 
// all of its functionalities, during different
// stages in the program, data will be passed
// between states which changes the behavior
// of its underlying functionalities. For 
// example, during `Genesis` stage, any control
// packet besides `Connect` results in immediate 
// disconnection from the broker. After `Genesis`,
// data will be passed to `Online` state which is
// opposite of `Genesis` state ( `Connect` results
// in immediate termination ).
type ConnectionState interface {
	ConStateInterface
	SetClient(client ClientInterface)
	SetServer(server ServerInterface)
}

// BaseControlInterface is the interface 
// to conform to fulfilling requirements
// of internal management console.
type BaseControlInterface interface {
  Shutdown()
}

