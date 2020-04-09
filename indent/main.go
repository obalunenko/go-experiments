package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

type user struct {
	ID      int64         `json:"id"`
	Name    string        `json:"name"`
	Created time.Time     `json:"created_at"`
	Expire  time.Duration `json:"expire_after"`
	Pem     string        `json:"-"`
}

// go binary encoder
func toGOB64(u user) (string, error) {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(u); err != nil {
		return "", errors.Wrap(err, `failed gob Encode`)
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// go binary decoder
func fromGOB64(str string) (user, error) {
	by, err := base64.StdEncoding.DecodeString(str)

	if err != nil {
		return user{}, errors.Wrap(err, `failed base64 Decode`)
	}

	var (
		b bytes.Buffer
		u user
	)

	b.Write(by)

	if err := gob.NewDecoder(&b).Decode(&u); err != nil {
		return user{}, errors.Wrap(err, `failed gob Decode`)
	}
	return u, nil
}

func prettyPrint(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal")
	}

	var out bytes.Buffer

	if err := json.Indent(&out, b, "", "\t"); err != nil {
		return "", errors.Wrap(err, "failed to indent")
	}

	if _, err := out.WriteString("\n"); err != nil {
		return "", errors.Wrap(err, "failed to write string")
	}

	return out.String(), nil
}

func main() {
	user1 := user{
		ID:      1234567,
		Name:    "Test",
		Created: time.Now(),
		Expire:  time.Minute * 3,
	}

	b64, err := toGOB64(user1)
	errFatal(err)

	fmt.Println(b64)
	userD, err := fromGOB64(b64)
	errFatal(err)

	s, err := prettyPrint(userD)
	errFatal(err)

	fmt.Println(s)
}

func errFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
