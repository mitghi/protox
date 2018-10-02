package commands

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWriteArray(t *testing.T) {
	var (
		c *Command = NewCommand()
	)
	c.WriteArrayHeader(2, 2)
	if bytes.Compare(c.Bytes(), []byte("*2:2\r\n")) != 0 {
		t.Fatal("inconsistent state.")
	}
	c.WriteArrayHeader(2, 2)
	c.WriteArrayItem("test")
	c.WriteArrayItem("test2")
	c.WriteArrayHeader(2, 2)
	c.WriteArrayItem("test3")
	c.WriteArrayItem("test4")
	fmt.Printf("%s\n", string(c.Bytes()))
	c.Reset()
	if c.size != 0 || c.valcnt != 0 || c.totalcnt != 0 {
		t.Fatal("inconsistent state.")
	}
	c.WriteArrayHeader(3, 12)
	c.WriteArrayItem("SET")
	c.WriteArrayItem("test")
	c.WriteArrayItem("test2")
	fmt.Println("")
	fmt.Printf("%s\n", string(c.Bytes()))
	if c.totalcnt != 12 {
		t.Fatal("assertion failed.")
	}
}

func TestMakeAndWriteArray(t *testing.T) {
	var (
		c  *Command = NewCommand()
		ac *Command // array command
	)
	ac = c.MakeArray()
	if ac == nil {
		t.Fatal("inconsistent state.")
	}
	ac.WriteArrayItem("SET")
	ac.WriteArrayItem("test")
	ac.WriteArrayItem("test2")
	if !ac.CommitArray() {
		t.Fatal("inconsistent state.")
	}
	fmt.Println(c.String())
}
