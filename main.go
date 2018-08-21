package main

import (
    "fmt"
    "github.com/ProtonMail/go-mbox"
    "encoding/json"
    "io"
    "os"
)

type Mbox struct {
    Messages []*Message
}

func Read(r io.Reader) (*Mbox, error) {
    var messages []*Message

    msgs := mbox.NewScanner(r)
    buf := make([]byte, 0, 64*1024)
    msgs.Buffer(buf, 1024*1024*100)
    for msgs.Next() {
        messages = append(messages, Decode(msgs.Message()))
    }

    return &Mbox{
        Messages: messages,
    }, msgs.Err()
}

func ReadFile(filename string) (*Mbox, error) {
    fp, err := os.Open(filename)
    if err != nil {
        return nil, err
    }

    return Read(fp)
}

func main() {
    for n, arg := range os.Args {
        if n == 0 {
            continue
        }

        if mbox, err := ReadFile(arg); err == nil {
            for _, mail := range mbox.Messages {

                j, err := json.Marshal(mail)
                if err != nil {
                    fmt.Println(err)
                    return
                }
                fmt.Println(string(j))
            }
        } else {
            fmt.Printf("%s\n", err.Error())
        }
    }
}