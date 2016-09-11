package scp

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"golang.org/x/crypto/ssh"
)

func SCPSERVER() {
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

// only create messages are supported for now. scp can preserve file attributes
// but that doesnt really make sense for us.
func parseHeader(header string) (int, string, error) {
	// C0644 4 test
	perm := ""
	size := 0
	name := ""
	_, err := fmt.Sscanf(header, "%s %d %s", &perm, &size, &name)

	if err != nil {
		return 0, "", err
	}

	return size, name, nil
}

// parses the ssh exec payload arguments. since only scp is supported the only
// things that are important are the file name and the command scp.
func parsePayload(payload []byte) (string, error) {
	// XXXXscp -t foobar
	// first four bytes are the length of the command.
	p := string(payload[4:])
	parts := strings.Split(p, " ")
	cmd := strings.Trim(parts[0], " ")
	fileName := strings.Trim(parts[len(parts)-1], " ")

	if cmd != "scp" {
		spew.Dump(cmd)
		return "", errors.New("Only scp is supported")
	}

	return fileName, nil
}

func handleChannel(ch ssh.NewChannel) {
	fmt.Println("Handling Channel: ", ch.ChannelType())
	if ch.ChannelType() != "session" {
		ch.Reject(ssh.UnknownChannelType, "unknown channel type")
		return
	}

	connection, requests, err := ch.Accept()
	if err != nil {
		fmt.Println("Can not accept channel:", err)
		return
	}

	log.Print("Creating File to write to...")
	x, err := os.OpenFile("foobartest", os.O_CREATE|os.O_WRONLY, 0777)

	close := func() {
		connection.Close()
		x.Close()
		fmt.Println("closeed")
	}

	if err != nil {
		log.Printf("Could not opne file (%s)", err)
		close()
		return
	}

	go func() {
		defer connection.Close()

		for req := range requests {
			spew.Dump(req.Type)
			switch req.Type {
			case "exec":
				dstName, err := parsePayload(req.Payload)

				if err != nil {
					fmt.Println(err)
					fmt.Fprintf(connection, err.Error())
					return
				}

				// Acknowledge payload.
				n, err := connection.Write([]byte{0})
				if err != nil {
					fmt.Println("error writing bytes", err)
				}
				fmt.Println("wrote", n, "bytes to connection")

				// Receive SCP Header
				scpHeader := make([]byte, 2048) // size of buf in openssh
				n, err = connection.Read(scpHeader)
				if err != nil {
					fmt.Println("Failed to read header", err)
					return
				}

				size, originalName, err := parseHeader(string(scpHeader))
				if err != nil {
					fmt.Println("Failed to parse header")
					fmt.Fprintln(connection, "Failed to understand your scp command")
					return
				}

				if len(dstName) == 0 {
					dstName = originalName
				}

				if len(dstName) < 3 {
					fmt.Fprintln(connection, "filename:", dstName, "is too short (minimum length = 3)")
					return
				}

				if size > 5*1024*1024 {
					fmt.Fprintf(connection, "File is over 1000 bytes\n")
					return
				}

				fmt.Println("Creating a file with the name", dstName, "with size", size)

				// Acknowledge We have received the SCP Header
				n, err = connection.Write([]byte{0})
				if err != nil {
					fmt.Println("Failed to reply", err)
					return
				}
				fmt.Println("Wrote", n, "bytes")

				// Read content of file
				fileDataBuf := make([]byte, size)
				n, err = io.ReadFull(connection, fileDataBuf)
				if err != nil {
					fmt.Println("Failed to read", err)
					return
				}
				fmt.Println("Read", n, "Bytes")

				// Set exit status
				connection.SendRequest("exit-status", false, []byte{0, 0, 0, 0})

				if err != nil {
					fmt.Println("reply error", err)
					return
				}

				// at this point we're good.
				fmt.Fprintf(connection, "\x00")
				err = req.Reply(true, nil)

				// Write the data to a file
				n, err = x.Write(fileDataBuf)
				if err != nil {
					fmt.Println("error writing bytes", err)
				}
				fmt.Println("Wrote", n, "bytes to disk")
				return
			}
		}
	}()
}
