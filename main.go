package main

import (
    "fmt"
    "github.com/blabber/mbox"
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

                b, err := json.Marshal(mail)
                if err != nil {
                    fmt.Println(err)
                    return
                }
                fmt.Println(string(b))

                // for k, vs := range mail.Header {
                //     for _, v := range vs {
                //         fmt.Printf("%s: %s\n", k, v)
                //     }
                // }
                // for _, body := range mail.Bodies {
                //     fmt.Println("====================================================")
                //     for k, vs := range body.Header {
                //         for _, v := range vs {
                //             fmt.Printf("%s: %s\n", k, v)
                //         }
                //     }
                //     fmt.Println("")
                //     io.Copy(os.Stdout, body.Content)
                //     fmt.Println("")
                // }
                // fmt.Println("====================================================\n\n")
            }
        } else {
            fmt.Printf("%s\n", err.Error())
        }
    }
}