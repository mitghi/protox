package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitghi/protox/auth"
	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/messages"
	"github.com/mitghi/protox/networking"
	"github.com/mitghi/protox/protobase"
	//	"github.com/pkg/profile"
)

var (
	logger *logging.Logging

	username    *string
	password    *string
	certificate *string
	key         *string
	sr          *string
	addr        string
	qos         *int
	user        *User
)

func init() {
	logger = logging.NewLogger("CLI")
}

func pcallback(opts protobase.OptionInterface, msg protobase.MsgInterface) {
	// TODO
	// . comment
}

func main() {
	//	defer profile.Start(profile.CPUProfile).Stop()
	log.Println("[+] started")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL)

	username = flag.String("username", "", "username")
	password = flag.String("password", "", "password")
	certificate = flag.String("certificate", "./client/client.pem", "")
	key = flag.String("key", "./client/client.key", "")
	qos = flag.Int("qos", 0, "quality of service")
	sr = flag.String("sr", "a/simple/route", "subscribe route")

	flag.Parse()
	if raddr := os.Getenv("PROTOX_ADDR"); raddr != "" {
		addr = raddr
	} else {
		addr = ":52909"
	}
	opts := client.CLBOptions{
		Addr:      addr,
		MaxRetry:  10,
		HeartBeat: 4,
		ClientDelegate: func() protobase.ClientInterface {
			cl := NewCustomClient(*username, *password, "")
			cl.SetCreds(&auth.Creds{*username, *password, ""})
			return cl
		},
		StorageDelegate: messages.NewMessageBox(),
		Conn: func(addr string) (clbc *networking.CLBConnection) {
			clbc = networking.NewClientConnection(addr)
			err := clbc.SetupTLSConfig(*certificate, *key)
			if err != nil {
				logger.Fatal("- [NewClientConnection] unable to setup TLS. err:", err)
				panic("unable to setup tls.")
			}
			return clbc
		}(addr),
		SecMRS:     5,
		CFCallback: nil,
	}
	ui, _ := NewCLUI()
	if ui == nil {
		panic("ui == nil")
	}
	user = NewUser(opts)
	if user == nil {
		panic("unable to create user")
	}
	cl := user.Cl.(*CustomClient)
	cl.user = user
	cl.ui = ui
	cl.ui.callback = cl.pubCli
	if err := user.Setup(); err != nil {
		panic("unable to setup user")
	}
	user.Connect()
	go ui.run()

	<-sigs
	user.Disconnect()
	logger.Debug("received signal, exiting....")
	<-user.Exch
}
