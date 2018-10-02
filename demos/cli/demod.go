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
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
//	"github.com/pkg/profile"
)

var (
	logger *logging.Logging

	username *string
	password *string
	sr       *string
	addr     string
	qos      *int
	user     *User
)

func init() {
	logger = logging.NewLogger("DemoD")
}

func pcallback(opts protobase.OptionInterface, msg protobase.MsgInterface) {

}

func main() {
//	defer profile.Start(profile.CPUProfile).Stop()
	log.Println("[+] started")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL)

	username = flag.String("username", "", "username")
	password = flag.String("password", "", "password")
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
		Conn:            protocol.NewClientConnection(addr),
		SecMRS:          5,
		CFCallback:      nil,
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
	log.Println("received signal, exiting....")
	<-user.Exch
}
