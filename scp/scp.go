package scp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/davecgh/go-spew/spew"

	"golang.org/x/crypto/ssh"
)

// type scpServer struct {
// 	ssh     easyssh.Server
// 	keyPath string
// }

// func NewSCPServer() *scpServer {
// 	s := easyssh.Server{Addr: ":2022"}

// 	privateBytes, err := ioutil.ReadFile("keys/id_rsa")
// 	if err != nil {
// 		fmt.Println("failed to load private key")
// 	}
// 	private, err := ssh.ParsePrivateKey(privateBytes)
// 	if err != nil {
// 		fmt.Println("failed to parse private key")
// 	}
// 	config := &ssh.ServerConfig{
// 		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
// 			return nil, nil
// 		},
// 	}

// 	config.AddHostKey(private)
// 	s.Config = config

// 	handler := easyssh.NewSessionServerHandler()
// 	channelHandler := easyssh.NewChannelsMux()

// 	channelHandler.HandleChannel(easyssh.SessionRequest, SessionHandler())

// 	handler.MultipleChannelsHandler = channelHandler
// 	s.Handler = handler

// 	x := &scpServer{
// 		ssh: s,
// 	}

// 	return x
// }

// func (s *scpServer) ListenAndServe() {
// 	s.ssh.ListenAndServe()
// }

func SCPSERVER() {

	//https://gist.githubusercontent.com/jedy/3357393/raw/e8e671080a8a04964d1a352fda167777aa163f1f/go_scp.go

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	privateBytes, err := ioutil.ReadFile("/home/dylan/.ssh/id_rsa")
	if err != nil {
		panic("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		panic("failed to listen for connection")
	}

	for {
		fmt.Println("Accepting")
		nConn, err := listener.Accept()
		if err != nil {
			panic("failed to accept incoming connection")
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		_, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			panic("failed to handshake")
		}
		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)
		// Accept all channels
		go handleChannels(chans)
	}
}

func handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go handleChannel(newChannel)
	}
}

func handleChannel(newChannel ssh.NewChannel) {
	fmt.Println("Handling Channel: ", newChannel.ChannelType())
	if newChannel.ChannelType() != "session" {
		newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
		return
	}

	connection, requests, err := newChannel.Accept()
	if err != nil {
		fmt.Println("Can not accept channel:", err)
		return
	}

	// 	bash := exec.Command("scp -t scptest")
	// 	close := func() {
	// 		connection.Close()
	// 		if err != nil {
	// 			log.Printf("Failed to exit scp (%s)", err)
	// 		}
	// 		log.Println("session closed")
	// 	}

	log.Print("Creating File to write to...")
	x, err := os.OpenFile("foobartest", os.O_CREATE|os.O_WRONLY, 0777)

	close := func() {
		connection.Close()
		x.Close()
		fmt.Println("closeed")
	}

	fmt.Println("extra: ")
	spew.Dump(newChannel.ExtraData())

	if err != nil {
		log.Printf("Could not opne file (%s)", err)
		close()
		return
	}

	// connection.Write(
	// bashf, err := pty.Start(bash)
	// if err != nil {
	// 	log.Printf("Could not start pty (%s)", err)
	// 	close()
	// 	return
	// }

	// var b bytes.Buffer

	// var once sync.Once
	// go func() {
	// 	io.Copy(connection, x)
	// 	once.Do(close)
	// }()
	// go func() {
	// 	io.Copy(x, connection)
	// 	once.Do(close)
	// }()

	//o. defer connection.Close()

	go func() {
		for req := range requests {
			spew.Dump(req.Type)
			switch req.Type {
			case "exec":
				fmt.Println("payload", req.Payload)
				fmt.Println("want reply", req.WantReply)

				// buf := make([]byte, 256)

				// scanner := bufio.NewScanner(connection)

				// for {

				n, err := connection.Write([]byte{0})
				fmt.Println("wrote", n, "bytes to connection")

				if err != nil {
					fmt.Println("error writing bytes", err)
				}

				buf := make([]byte, 256)
				n, err = connection.Read(buf)

				if err != nil {
					fmt.Println("Failed to receive from connection", err)
					return
				}

				fmt.Println("buf:")
				spew.Dump(buf)

				n, err = connection.Write([]byte{0})
				if err != nil {
					fmt.Println("Failed to reply", err)
					return
				}

				fmt.Println("Wrote", n, "bytes")

				buf = make([]byte, 10)

				fmt.Println("b4 read")
				n, err = connection.Read(buf)
				fmt.Println("after read")

				if err != nil {
					fmt.Println("Failed to read", err)
					return
				}

				fmt.Println("Read", n, "Bytes")

				spew.Dump(buf)

				// }
				err = req.Reply(true, []byte(""))
				if err != nil {
					fmt.Println("reply err", err)
				}

				// connection.Write([]byte("\n\n"))

				// connection.Read(buf)
				// spew.Dump("buf")
				// spew.Dump(buf)
				// // fmt.Println("before peek")
				// spew.Dump(newChannel.ExtraData())
				// stuff := bufio.NewReader(newChannel)
				// bitz, err := stuff.Peek(5)
				// fmt.Println("after peek")

				// if err != nil {
				// 	fmt.Println("failed to peek", err)
				// }

				// fmt.Println("bitz: ", bitz)
			}
		}
	}()
}
