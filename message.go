package main

import (
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
}

type Body struct {
    ContentType string `json:"Content-Type"`
    Content     string
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

func newFileByPart(part *multipart.Part) *File {
    headers := make(map[string][]string)
    file := &File{
        FileName: "",
        ContentType: "",
        Content:  nil,
    }

    for k, v := range part.Header {
        headers[http.CanonicalHeaderKey(k)] = decodeHeaders(v)
    }

    file.FileName = decodeHeader(part.FileName())

    buf := bytes.NewBuffer([]byte{})
    _, err := io.Copy(buf, part)
    if err != nil {
        return nil
    }

    mediaType, params, err := mime.ParseMediaType(headers["Content-Type"][0])
    if err != nil {
        return nil
    }

    charset := params["charset"]
    encoding := headers["Content-Transfer-Encoding"][0]

    file.ContentType = mediaType

    file.Content = newDecoder(buf, charset, encoding)

    return file
}

func newBodyByMessage(message *Message, header string, r io.Reader) *Body {
    body := &Body{
        ContentType:     "",
        Content:         "",
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

    body.Content = buf.String()

    return body
}