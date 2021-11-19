package satel

import (
	"bufio"
	"errors"
	"net"
	"sync"
	"time"
)

type Event struct {
	Type  ChangeType
	Index int
	Value bool
}

type Config struct {
	EventsQueueSize int
	LongCommands    bool
}

type Satel struct {
	conn    net.Conn
	mu      sync.Mutex
	cmdSize int
	cmdChan chan int
	Events  chan Event
}

func New(conn net.Conn) *Satel {
	return NewConfig(conn, Config{})
}

func NewConfig(conn net.Conn, config Config) *Satel {
	s := &Satel{
		conn:    conn,
		cmdChan: make(chan int),
		Events:  make(chan Event, config.EventsQueueSize),
	}
	if config.LongCommands {
		s.cmdSize = 32
	} else {
		s.cmdSize = 16
	}
	go s.read()
	err := s.sendCmd(0x7F, 0x01, 0x04, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
	if err != nil {
		close(s.Events)
		return s
	}
	go func() {
		for {
			err = s.sendCmd(0x0A)
			if err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return s
}

func (s *Satel) ArmPartition(code string, mode, index int) error {
	data := make([]byte, 4)
	data[index/8] = 1 << (index % 8)
	bytes := prepareCommand(code, byte(0x80+mode), data...)
	return s.sendCmd(bytes...)
}

func (s *Satel) ForceArmPartition(code string, mode, index int) error {
	data := make([]byte, 4)
	data[index/8] = 1 << (index % 8)
	bytes := prepareCommand(code, byte(0xA0+mode), data...)
	return s.sendCmd(bytes...)
}

func (s *Satel) DisarmPartition(code string, index int) error {
	data := make([]byte, 4)
	data[index/8] = 1 << (index % 8)
	bytes := prepareCommand(code, byte(0x84), data...)
	return s.sendCmd(bytes...)
}

func (s *Satel) SetOutput(code string, index int, value bool) error {
	cmd := byte(0x89)
	if value {
		cmd = 0x88
	}
	data := make([]byte, s.cmdSize)
	data[index/8] = 1 << (index % 8)
	bytes := prepareCommand(code, cmd, data...)
	return s.sendCmd(bytes...)
}

func prepareCommand(code string, cmd byte, data ...byte) []byte {
	bytes := append([]byte{cmd}, transformCode(code)...)
	return append(bytes, data...)
}

func (s *Satel) Close() error {
	return s.conn.Close()
}

type command struct {
	prev        [32]byte
	initialized bool
}

func (s *Satel) read() {
	scanner := bufio.NewScanner(s.conn)
	scanner.Split(scan)
	commands := make(map[byte]command)

	for ok := scanner.Scan(); ok; ok = scanner.Scan() {
		bytes := scanner.Bytes()
		cmd := bytes[0]
		bytes = bytes[1 : len(bytes)-2]
		s.cmdRes()
		if cmd == 0xEF {
			continue
		}
		c := commands[cmd]
		for i, bb := range bytes {
			change := bb ^ c.prev[i]
			for j := 0; j < 8; j++ {
				index := byte(1 << j)
				if !c.initialized || change&index != 0 {
					s.Events <- Event{
						Type:  ChangeType(cmd),
						Index: i*8 + j,
						Value: bb&index != 0,
					}
				}
			}
			c.prev[i] = bytes[i]
		}
		c.initialized = true
		commands[cmd] = c
	}
	close(s.Events)
	_ = s.conn.Close()
}

func (s *Satel) cmdRes() {
	select {
	case s.cmdChan <- 0:
	default:
	}
}

func (s *Satel) sendCmd(data ...byte) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return errors.New("no connection")
	}
	_, err = s.conn.Write(frame(data...))
	if err == nil {
		select {
		case <-s.cmdChan:
		case <-time.After(3 * time.Second):
		}
	}
	return
}

func transformCode(code string) []byte {
	bytes := make([]byte, 8)
	for i := 0; i < 16; i++ {
		if i < len(code) {
			digit := code[i]
			if i%2 == 0 {
				bytes[i/2] = (digit - '0') << 4
			} else {
				bytes[i/2] |= digit - '0'
			}
		} else if i%2 == 0 {
			bytes[i/2] = 0xFF
		} else if i == len(code) {
			bytes[i/2] |= 0x0F
		}
	}
	return bytes
}
