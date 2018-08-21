package main

import (
    "encoding/hex"
    "crypto/sha256"
    "bytes"
    "io"
    "mime"
    "mime/multipart"
    "net/http"
)

type Message struct {
    Header Header
    Body   []*Body
    Files  []*File
}

type Header map[string][]string

type File struct {
    FileName    string
    ContentType string `json:"Content-Type"`
    Content     io.Reader
    Base64      string `json:"base64"`
    Sha256      string `json:"sha256"`
}

type Body struct {
    ContentType string `json:"Content-Type"`
    Text        string
}

func (h Header) Get(key string) string {
    if s, ok := h[http.CanonicalHeaderKey(key)]; ok && len(s) > 0 {
        return s[0]
    }
    return ""
}

func (h Header) Set(key, val string) {
    h[http.CanonicalHeaderKey(key)] = []string{val}
}

func (h Header) Del(key string) {
    delete(h, http.CanonicalHeaderKey(key))
}

func newFile(part *multipart.Part) *File {
    headers := make(map[string][]string)
    file := &File{
        FileName:    "",
        ContentType: "",
        Content:     nil,
        Base64:      "",
        Sha256:      "",
    }
    hash := sha256.New()

    for k, v := range part.Header {
        headers[http.CanonicalHeaderKey(k)] = decodeHeaders(v)
    }

    mediaType, params, err := mime.ParseMediaType(headers["Content-Type"][0])
    if err != nil {
        return nil
    }

    charset := params["charset"]
    encoding := headers["Content-Transfer-Encoding"][0]

    file.ContentType = mediaType

    file.FileName = decodeHeader(part.FileName())

    buf1 := bytes.NewBuffer([]byte{})
    buf2 := bytes.NewBuffer([]byte{})

    writer := io.MultiWriter(hash, buf1, buf2)
    if _, err := io.Copy(writer, part); err != nil {
        return nil
    }

    file.Content = newDecoder(buf1, charset, encoding)

    file.Base64 = base64Encode(newDecoder(buf2, charset, encoding))

    file.Sha256 = hex.EncodeToString(hash.Sum(nil))

    return file
}

func newBody(message *Message, header string, r io.Reader) *Body {
    body := &Body{
        ContentType: "",
        Text:        "",
    }

    mediaType, params, err := mime.ParseMediaType(header)
    if err != nil {
        return nil
    }

    charset := params["charset"]
    encoding := message.Header.Get("Content-Transfer-Encoding")

    body.ContentType = mediaType

    buf := new(bytes.Buffer)
    buf.ReadFrom(newDecoder(r, charset, encoding))

    body.Text = buf.String()

    return body
}