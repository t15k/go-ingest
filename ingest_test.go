package ingest

import (
	"testing"
)

func TestBootstrapWithDConfig(t *testing.T) {
	type Mod1 struct {
		name string
	}
	RegisterMod("testmod1", func(c string) interface{} {
		t.Log(c)
		return &Mod1{c}
	})
	RegisterMod("testmod2", func(c string) interface{} {
		return &Mod1{"aa"}
	})
	mods, err := Bootstrap("-testmod1,remote:4343.-testmod2.")
	if err != nil {
		t.Fatal(err)
	}
	mod1 := mods[0]
	mod1I := mod1.(*Mod1)

	if mod1I.name != "remote:4343" {
		t.Fatalf("mod1 name was not remote:4343 but %s",
			mod1I.name)
	}
	mod2 := mods[1]
	mod2I := mod2.(*Mod1)
	if mod2I.name != "aa" {
		t.Fatal("mod2 name was not aa")
	}
	// mp>socko,remote:4343.>mpo..")
	// ">mp>socko,remote:4343.>mpo.."
}

type Mod1 struct {
	name string
	recs []Receiver
}

func (m *Mod1) AddReciever(r Receiver) {
	if m.recs == nil {
		m.recs = make([]Receiver, 0, 5)
	}
	m.recs = append(m.recs, r)
}

type TestReceiver struct {
	name string
}

func (t *TestReceiver) Accept(s interface{}) {

}

func TestBootstrapWithNestedConfig(t *testing.T) {
	RegisterMod("testmod3", func(c string) interface{} {
		t.Log(c)
		return &TestReceiver{c}
	})
	RegisterMod("testmod2", func(c string) interface{} {
		return &TestReceiver{"22"}
	})
	RegisterMod("testmod1", func(c string) interface{} {
		return &Mod1{"11", nil}
	})
	mods, err := Bootstrap("-testmod1-testmod2.-testmod3,host..")
	if err != nil {
		t.Fatal(err)
	}
	mod1 := mods[0]
	mod1I := mod1.(*Mod1)

	if mod1I.name != "11" {
		t.Fatalf("mod1 name was not 11 but %s",
			mod1I.name)
	}
	if mod1I.recs[0].(*TestReceiver).name != "22" {
		t.Fatal("Not 22")
	}
	if mod1I.recs[1].(*TestReceiver).name != "11" {
		t.Fatal("Not 11")
	}

}
