package scp

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/dependencies"

	"golang.org/x/crypto/ssh"
)

var (
	zeroByte = []byte{0}
	scpCmd   = []byte("scp")

	bindAddr = ":2020"
)

// Server foo
type Server struct {
	dependencies.Dependencies

	CertPath string
	BindAddr string
}

// ListenAndServe starts a SCP server.
func (s *Server) ListenAndServe() {
	if len(s.BindAddr) == 0 {
		panic("Must specify BindAddr")
	}
	if len(s.CertPath) == 0 {
		panic("Must specify CertPath")
	}

	privateBytes, err := ioutil.ReadFile(s.CertPath)
	if err != nil {
		panic("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", s.BindAddr)
	if err != nil {
		s.Fatal("Failed to bind to address", "addr", bindAddr)
	}

	s.Info("Listening for SCP connections", "addr", bindAddr)

	for {
		nConn, err := listener.Accept()
		if err != nil {
			s.Error("Failed to listen for SCP connections", "err", err)
			continue
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		_, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			s.Error("Failed to create SCP connection", "err", err)
			continue
		}

		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)

		// Accept all channels
		go s.handleChannels(chans)
	}
}

func (s *Server) handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go s.handleChannel(newChannel)
	}
}

// Only create (CXXXX) messages are supported for now. SCP can preserve file
// attributes but that doesn't make sense for zqz.
// Format: C0644 4 test
func parseHeader(header []byte) (perm string, size int64, name string, err error) {
	frags := strings.Split(string(header), " ")

	if len(frags) != 3 {
		err = errors.New("failed to parse header, not 3 fragments")
		return
	}

	if size, err = strconv.ParseInt(frags[1], 10, 64); err != nil {
		return
	}

	perm = frags[0]
	name = frags[2]

	return
}

// parses the ssh exec payload arguments. since only scp is supported the only
// things that are important are the file name and the command scp.
// XXXXscp -t foobar
// first four bytes are the length of the command.
func parsePayload(payload []byte) (string, error) {
	length := binary.BigEndian.Uint32(payload[:4])

	if length > 1024 {
		return "", errors.Errorf("Command too long. Max length: 1024 chars")
	}

	p := payload[4 : 4+length]

	parts := bytes.Split(p, []byte{' '})
	cmd := bytes.TrimSpace(parts[0])
	fileName := string(bytes.TrimSpace(parts[len(parts)-1]))

	if bytes.Compare(cmd, scpCmd) != 0 {
		return "", errors.Errorf("non-scp command was used: %q", p)
	}

	return fileName, nil
}

func (s *Server) handleChannel(ch ssh.NewChannel) {
	id := rand.Int()
	s.Debug("Handling Channel", "id", id, "chan", ch.ChannelType())

	if ch.ChannelType() != "session" {
		s.Info("Received unknown channel type", "chan", ch.ChannelType())
		ch.Reject(ssh.UnknownChannelType, "unknown channel type")
		return
	}

	channel, requests, err := ch.Accept()
	if err != nil {
		s.Error("Failed to accept channe", "err", err)
		return
	}

	var closer sync.Once
	closeChannel := func() {
		s.Debug("Closed Channel", "id", id)
		channel.Close()
	}

	defer closer.Do(closeChannel)

	for req := range requests {
		spew.Dump(req.Type)
		switch req.Type {
		case "exec":
			// Let it through
		case "env":
			if req.WantReply {
				if err = req.Reply(true, nil); err != nil {
					s.Error("Failed to ignore env command", "err", err)
				}
			}
			continue
		default:
			s.Info("Received unhandled request type", "type", req.Type)
			continue
		}

		r := &scpRequest{db: s.DB}
		processors := []processor{
			r.ParseSCPRequest, r.DownloadFile, r.EndConnectionGracefully,
		}

		for _, proc := range processors {
			if err := proc(channel, req); err != nil {
				fmt.Fprintln(channel, "failed to process request:", err.Error())
				// log.Printf("%+v", err)
				break
			}
		}

		closer.Do(closeChannel)
	}
}

type processor func(channel ssh.Channel, req *ssh.Request) error

type scpRequest struct {
	size     int64
	original string
	filename string
	hash     string
	db       db.Executor
}

func (s *scpRequest) ParseSCPRequest(channel ssh.Channel, req *ssh.Request) error {
	var err error

	// Parse the payload received from the scp client
	if s.original, err = parsePayload(req.Payload); err != nil {
		return err
	}
	log.Println("being sent file:", s.original)

	// Acknowledge payload.
	if _, err = channel.Write(zeroByte); err != nil {
		return errors.Wrap(err, "failed to write")
	}

	// Receive SCP Header
	scpHeader := make([]byte, 2048) // size of buf in openssh
	if _, err = channel.Read(scpHeader); err != nil {
		return errors.Wrap(err, "failed to retrieve header")
	}

	if _, s.size, _, err = parseHeader(scpHeader); err != nil {
		return errors.Wrap(err, "failed to parse scp header")
	}

	// Acknowledge We have received the SCP Header
	if _, err = channel.Write(zeroByte); err != nil {
		return errors.Wrap(err, "failed to reply to scp header")
	}

	return nil
}

func (s *scpRequest) DownloadFile(channel ssh.Channel, req *ssh.Request) error {
	uid := uuid.NewV4()
	s.filename = filepath.Join(os.TempDir(), uid.String())

	log.Println("saving to:", s.filename)
	f, err := os.OpenFile(s.filename, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return errors.Wrap(err, "failed to open file for downloading")
	}

	h := sha1.New()
	mw := io.MultiWriter(f, h)
	// Read file contents
	var n int64
	if n, err = io.CopyN(mw, channel, s.size); err != nil {
		return errors.Wrap(err, "failed to download file")
	}
	log.Println("copied", n, "bytes")

	if err = f.Close(); err != nil {
		return errors.Wrap(err, "failed to close file")
	}

	s.hash = fmt.Sprintf("%x", h.Sum(nil))

	fmt.Println("Sha'd file", s.hash)

	return nil
}

func (s *scpRequest) EndConnectionGracefully(channel ssh.Channel, req *ssh.Request) error {
	ch := makeChanReq(channel, req)

	// Tell them we received their file
	ch.Write(zeroByte)
	// Set exit status of the "exec'd" program
	ch.SendRequest("exit-status", []byte{0, 0, 0, 0})

	return ch.err
}

type chanReq struct {
	channel ssh.Channel
	req     *ssh.Request

	err error
}

func makeChanReq(channel ssh.Channel, req *ssh.Request) *chanReq {
	return &chanReq{
		channel: channel,
		req:     req,
	}
}

func (c *chanReq) Write(b []byte) {
	if c.err != nil {
		return
	}

	n, err := c.channel.Write(b)
	if n != len(b) {
		c.err = errors.New("short write")
	}

	if err != nil {
		c.err = errors.Wrap(err, "failed to write to channel")
	}
}

func (c *chanReq) SendRequest(name string, payload []byte) {
	if c.err != nil {
		return
	}

	if _, err := c.channel.SendRequest(name, false, payload); err != nil {
		c.err = errors.Wrap(err, "failed to send request")
	}
}

func (c *chanReq) Reply(ok bool, payload []byte) {
	if c.err != nil {
		return
	}

	if err := c.req.Reply(ok, payload); err != nil {
		c.err = errors.Wrap(err, "failed to send reply")
	}
}
