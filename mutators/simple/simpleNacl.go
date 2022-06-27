//go:build crappy
// +build crappy

// this is tagged "crappy" because I don't like my code, but the nacl module is great.

package mutators

import (
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"github.com/batmac/ccat/log"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"
)

func init() {
	// we want the output to be as-is
	simpleRegister("easyseal", easyseal, withDescription("encrypt with Nacl EasySeal, key used is printed on stderr"),
		withCategory("encrypt"),
		withExpectingBinary(true))
	simpleRegister("easyopen", easyopen, withDescription("decrypt with Nacl EasyOpen, get the key from env (KEY)"),
		withCategory("decrypt"))
}

func easyseal(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	u, err := ioutil.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	key := getKey()
	box := secretbox.EasySeal(u, key)
	_, err = w.Write(box)
	if err != nil {
		return 0, err
	}
	defer func() {
		log.Printf("KEY=%s", hex.EncodeToString((*key)[:]))
	}()
	return int64(len(box)), nil
}

func easyopen(w io.WriteCloser, r io.ReadCloser) (int64, error) {
	u, err := ioutil.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}
	key := getKey()
	box, err := secretbox.EasyOpen(u, key)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(box)
	if err != nil {
		return 0, err
	}

	err = r.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
	return int64(len(box)), nil
}

func getKey() nacl.Key {
	var key nacl.Key
	if keyString := os.Getenv("KEY"); len(keyString) == 0 {
		key = nacl.NewKey()
	} else {
		var err error
		key, err = nacl.Load(keyString)
		if err != nil {
			log.Fatal(err)
		}
	}

	return key
}
