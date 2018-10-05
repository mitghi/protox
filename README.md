# protox

Implementation of high performance [protox](https://protox.xyz) message broker.

version : **Alpha**.

- [X] Fully Binary
- [X] Publish/Subscribe
- [X] Quality of Service 0 and 1 (at most once, at least once)
- [X] Persistent states
- [X] Client Library

Whitebox test suits

- [X] containers
- [X] core
- [X] messages
- [X] server

# Objectives

objectives for **Beta** version:

- [X] Retain Storage
- [X] Permissions
- [ ] one-to-one  Request/Response
- [ ] Message Queue
- [ ] Config parser

**TODO:** whitebox test suits

- [X] auth
- [X] protocol
- [ ] broker
- [ ] client
- [ ] utils


# Desired features

This is a wishlist for future versions.

- [ ] [Raft Consensus Algorithm](https://raft.github.io/raft.pdf)
- [ ] Proposals-over-Network (ex. job delegation, polls, stable matching, .... )
- [ ] One-to-Many Request/Response
- [ ] Buffered   Channels
- [ ] Unbuffered Channels
- [ ] Event Notifications
- [ ] Event Multiplexing
- [ ] Endpoints ( compatibility with 3rd party services )
- [ ] Management Console
- [ ] Signed Messages ( RSA, PGP, SHA256, SCRYPT)
- [ ] Pluggable Authenication subsystem
- [ ] Plugins subsystem

# Setup ( development )

[![](./media/deploy.gif)](https://asciinema.org/a/70y88r89i5pmqg4n9er8gb2oy)

using Docker
------------

**NOTICE:** create deploy key and save it as `id_rsa_depl` in `docker/`. Check `docker logs -f [CONTAINER ID]` for logs. For more options check `docker/Makefile`. Compiled binary is located at `$GOPATH/bin`.

to **build**

```bash
$ git clone git@github.com:mitghi/protox.git .
$ export PROTOXREPO="$GOPATH/src/github.com/mitghi"
$ mkdir -p $PROTOXREPO
$ mv ./protox $PROTOXREPO/protox && cd $PROTOXREPO/protox/docker
$ make build-remote
$ make shell
```

to **deploy**

```bash
$ git clone git@github.com:mitghi/protox.git .
$ export PROTOXREPO="$GOPATH/src/github.com/mitghi"
$ mkdir -p $PROTOXREPO
$ mv ./protox $PROTOXREPO/protox && cd $PROTOXREPO/protox/docker
$ make build
$ export PORTS="-p 52909:52909"
$ make spawn
```

After starting the container, use following snippet to send a crafted `Connect` packet to assure that broker can receive connections from outside.

```bash
$ echo -e '\x10''\x4f''\x0''\x4''\x50''\x52''\x58''\x31''\xf''\x0''\x3c''\x24''\x32''\x61''\x24''\x30''\x36''\x24''\x6c''\x4e''\x69''\x38''\x48''\x35''\x6b''\x63''\x35''\x5a''\x39''\x54''\x39''\x78''\x4a''\x41''\x58''\x77''\x51''\x71''\x79''\x75''\x6e''\x6c''\x32''\x45''\x59''\x68''\x47''\x55''\x69''\x36''\x63''\x74''\x33''\x54''\x67''\x70''\x52''\x31''\x42''\x4e''\x62''\x31''\x76''\x70''\x7a''\x70''\x70''\x39''\x70''\x7a''\x43''\x0''\x0''\x0''\x1''\x0''\x4''\x74''\x65''\x73''\x74' | nc localhost $(printf "%d" 0xcead) | hexdump
```

**(Client)Output:**

```
0000000 20 01 00
0000003
```

**(Broker)Output:**

```
$ RLOG_LOG_LEVEL=DEBUG protox
2017/03/14 13:41:46 [+]Started
2017-03-14T13:41:52+01:00 INFO     : * [Genesis] Participation request accepted. <- (*Server).ServeTCP
2017-03-14T13:41:52+01:00 DEBUG    : * [Packet] raw packet content 0x10 0x4f 0x00 0x04 0x50 0x52 0x58 0x31 0x0f 0x00 0x3c 0x24 0x32 0x61 0x24 0x30 0x36 0x24 0x6c 0x4e 0x69 0x38 0x48 0x35 0x6b 0x63 0x35 0x5a 0x39 0x54 0x39 0x78 0x4a 0x41 0x58 0x77 0x51 0x71 0x79 0x75 0x6e 0x6c 0x32 0x45 0x59 0x68 0x47 0x55 0x69 0x36 0x63 0x74 0x33 0x54 0x67 0x70 0x52 0x31 0x42 0x4e 0x62 0x31 0x76 0x70 0x7a 0x70 0x70 0x39 0x70 0x7a 0x43 0x00 0x00 0x00 0x01 0x00 0x04 0x74 0x65 0x73 0x74 <- (*ProtoConnection).HandleDefault

2017-03-14T13:41:52+01:00 DEBUG    : (--OPTIONS[keepalive, clid, clusrname, clpasswd]=( true true true true )--)
2017-03-14T13:41:52+01:00 DEBUG    : * [Packet] conn packet content.  <- (*ProtoConnection).HandleDefault
	Username(test), Password($2a$06$lNi8H5kc5Z9T9xJAXwQqyunl2EYhGUi6ct3TgpR1BNb1vpzpp9pzC),
	ClientId(), KeepAlive(1), Version(PRX1)
2017-03-14T13:41:52+01:00 DEBUG    : * [clientDelegate] client [test] has joined.
2017-03-14T13:41:52+01:00 DEBUG    : * [Auth] status for uid [test] is [OK].
2017-03-14T13:41:52+01:00 DEBUG    : + [Genesis] for client [test] changed to [ready].
2017-03-14T13:41:52+01:00 DEBUG    : + [Client] Connected.
2017-03-14T13:41:52+01:00 DEBUG    : + [Client] Passed (Genesis) state and is now (Online). <- (*Server).NotifyDisconnected
2017-03-14T13:41:52+01:00 ERROR    : - [RecvHandler] error while receiving packets. <- (*ProtoConnection).recvHandler
2017-03-14T13:41:57+01:00 DEBUG    : - [Connection] DEADLINE, No heartbeat from client [test], terminating ....
2017-03-14T13:41:57+01:00 DEBUG    : - [1][Error] STATERR.
2017-03-14T13:41:57+01:00 DEBUG    : * [1][Event] Shutting down stream.
2017-03-14T13:41:57+01:00 DEBUG    : * [Connection] beginning to wait for coroutines to finish** <- (*ProtoConnection).**(%!s(MISSING))
2017-03-14T13:41:57+01:00 DEBUG    : - [2][Event][Death] Detached [Client]. userId test <- (*Server).NotifyDisconnected
2017-03-14T13:41:57+01:00 DEBUG    : - [2][Event][ConnEnded] for [Client]. userId test  <- (*Server).NotifyDisconnected
2017-03-14T13:41:57+01:00 DEBUG    : - [Client] connection to client [test] is terminated.
2017-03-14T13:41:57+01:00 DEBUG    : * [Connection] all coroutines are finished, exiting connection handler++ <- (*ProtoConnection).++(%!s(MISSING))
```


Native
------
For compiling the repository, latest version of [Golang](https://golang.org) is required.
Check Golang [installation guide](https://golang.org/doc/install) for futher details.
This is a general and generic way to compile Protox on **nix*.

```bash
$ mkdir -p $HOME/work && mkdir -p $HOME/work/bin && mkdir -p $HOME/work/src/github.com/mitghi
$ export $GOPATH="$HOME/work"
$ cd $HOME/work && pushd .;cd src/github.com/mitghi
$ git clone git@github.com:mitghi/protox.git ./protox && popd
$ go get ./... && go install -v github.com/mitghi/protox
$ cd ./bin && ./protox
```

# Example

Using sample Protox **Broker** with explicit **Access Control List**:


```go
package main

import (
	"os"

	"github.com/mitghi/protox/auth"
	"github.com/mitghi/protox/broker"
	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/protobase"
)

// clientDelegate is the delegate used by server to
// create a compatible `protobase.ClientInterface`
// struct.
func clientDelegate(uid string, pid string, cid string) protobase.ClientInterface {
	ret := User{client.Client{Username: uid, Password: pid, ClientId: cid}}
	return &ret
}

// defaultAuthConfig returns default configuration
// of Auth subsystem ( Strict Mode ) with a few
// dummy clients.
func defaultAuthConfig() *auth.AuthConfig {
	// initialize auth configuration container
	var c *auth.AuthConfig = auth.NewAuthConfig()
	// set authorization mode
	c.Mode = protobase.AUTHModeStrict
	// set access group partitions
	c.AccessGroups = auth.AuthGroups{
		Members: map[string][][3]string{
			"User": [][3]string{
				{"can", "publish", "self/inbox"},
				{"can", "subscribe", "self/notifications"},
			},
			"Bot": [][3]string{
				{"can", "publish", "self/location"},
				{"can", "request", "self/access/upgrade"},
			},
		},
		// set access control list authorization
		// type.
		Type: protobase.ACLModeInclusive,
	}
	// define sample credentials
	c.Credentials = []auth.AuthEntity{

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test",
				Password: "$2a$06$lNi8H5kc5Z9T9xJAXwQqyunl2EYhGUi6ct3TgpR1BNb1vpzpp9pzC",
				ClientId: "",
			},
			// define permission partition
			Group: "User",
		},

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test2",
				Password: "$2a$06$2uqusEvRMcpla2KXph8sBuBXO4WVOgIVbIgfRjk5y01UXxxgR9z6O",
				ClientId: "",
			},
			// define permission partition
			Group: "User",
		},

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test3",
				Password: "$2a$06$sgQ9yjjVvRxQhLqWKSGv4OTE2EF4ojUu1sEHnGUJdimmn.5M9M7/.",
				ClientId: "",
			},
			// define permission partition
			Group: "Bot",
		},

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test4",
				Password: "$2a$06$9wavlAtmNZ66Whe2wturDO7yIBdE41/Zcn4c5z4ydzJ/ydVJIZwJK",
				ClientId: "",
			},
			// define permission partition
			Group: "Bot",
		},
	}

	return c
}

var (
	c       *auth.AuthConfig        = defaultAuthConfig() // authentication configs
	authsys protobase.AuthInterface                       // authentication subsystem
)

func main() {
	var err error
	// setup authentication subsystem from configurations
	authsys, err = auth.NewAuthenticatorFromConfig(c)
	if err != nil {
		panic(err)
	}
	// initialize the broker with configured
	// authentication subsystem.
	brk := broker.NewBroker(broker.Options{
		Auth: authsys,
	}).(*broker.Broker)
	// run the broker
	ok := brk.Start()
	if !ok {
		os.Exit(1)
	}
	// wait for termination conditions such
	// as (KILL SIGNAL, FATAL ERRORS, .... ).
	<-brk.E

	os.Exit(0)
}
```

Check more examples at [demos](./demos/).

# Client

List of protox client libraries.

- [Protox Client](https://github.com/mitghi/protox_client) ( Golang ).

