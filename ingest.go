package ingest

import (
	"errors"
	"fmt"
	"strings"
	"text/scanner"
	"unicode"
)

var mods = map[string]Factory{}

// Receiver should be attached to an input.
type Receiver interface {
	Receive(s interface{})
}

// Factory produces receiver and or Emitters from cfg.
type Factory func(cfg string) interface{}

// Emitter is implemented by handlers that will emit data.
type Emitter interface {
	AddReceiver(r Receiver)
}

// SimpleEmitter just holds receivers. Can be nested in more complex receivers.
type SimpleEmitter struct {
	receivers []Receiver
}

// AddReceiver to this.
func (se *SimpleEmitter) AddReceiver(r Receiver) {
	if se.receivers == nil {
		se.receivers = make([]Receiver, 0, 2)
	}
	se.receivers = append(se.receivers)
}

// Emit o to all receivers add to this.
func (se *SimpleEmitter) Emit(o interface{}) {
	for _, rec := range se.receivers {
		rec.Receive(o)
	}
}

// RegisterMod register a module so if can be found by configuration.
func RegisterMod(name string, fact Factory) {
	mods[name] = fact
}

// Bootstrap will return an app initialised from config.
func Bootstrap(config string) ([]interface{}, error) {
	var s scanner.Scanner
	s.IsIdentRune = isIdentRune
	s.Init(strings.NewReader(config))
	tok := s.Scan()
	mods := make([]interface{}, 0, 5)
	for tok != scanner.EOF {
		txt := s.TokenText()
		if txt != "-" {
			return nil, errors.New("Must start with -")
		}
		vv, err := handleArrow(&s)
		if err != nil {
			return nil, err
		}
		for _, v := range vv {
			mods = append(mods, v)
		}
		tok = s.Scan()
	}
	if len(mods) == 0 {
		return nil, errors.New("No config")
	}
	return mods, nil
}

func handleConfig(s *scanner.Scanner) (string, error) {
	token := s.Scan()
	if token != scanner.EOF {
		return s.TokenText(), nil
	}
	return "", errors.New("EOF")
}

func handleArrow(s *scanner.Scanner) (sr []interface{}, rerr error) {
	tok := s.Scan()
	if tok == scanner.EOF {
		rerr = errors.New("Missing identifier after -")
		return
	}
	name := s.TokenText()
	fact, ok := mods[name]
	if !ok {
		rerr = fmt.Errorf("No module named %v", name)
		return
	}
	sr = make([]interface{}, 0, 5)
	var (
		next         string
		cf           string
		subReceivers = make([]interface{}, 0, 5)
	)
	for next != "." && tok != scanner.EOF {
		tok = s.Scan()
		next = s.TokenText()
		switch next {
		case "-":
			subReceivers, rerr = handleArrow(s) // add as next
			if rerr != nil {
				return
			}
		case ",":
			cf, rerr = handleConfig(s) // init with config
			if rerr != nil {
				return
			}
		}
	}
	mod := fact(cf)
	emitter, ok := mod.(Emitter)
	if ok {
		for _, rec := range subReceivers {
			rec, ok := rec.(Receiver)
			if ok {
				emitter.AddReceiver(rec)
			}
		}
	}
	sr = append(sr, mod)
	return
}

func isIdentRune(ch rune, i int) bool {
	switch ch {
	case '.':
		return i == 0
	}
	return ch == '_' || ch == ':' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
}
