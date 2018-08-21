package main

import (
    "bytes"
    "encoding/base64"
    "io"
    "mime"
    "mime/multipart"
    "mime/quotedprintable"
    "net/http"
    "net/mail"
    "regexp"
    "strings"
)

var encodedHeaderRegexp = regexp.MustCompile("\\=\\?([^?]+)\\?([BQ])\\?(.*?)\\?\\=\\s*")

func Decode(msg *mail.Message) *Message {
    message := &Message{
        Header:  Header{},
        XHeader: XHeader{},
        Body:    []*Body{},
        Files:   []*File{},
    }

    for key, header := range msg.Header {
        message.XHeader[http.CanonicalHeaderKey(key)] = decodeHeaders(header)
    }

    mediaType, params, err := mime.ParseMediaType(message.XHeader.Get("Content-Type"))
    if err != nil {
        return message
    }

    if strings.HasPrefix(mediaType, "multipart/") {
        message.XHeader.Set("Content-Type", mediaType)
        mr := multipart.NewReader(msg.Body, params["boundary"])
        for {
            p, err := mr.NextPart()
            if err == io.EOF {
                break
            }
            if err != nil {
                break
            }

            part := newFile(p)

            if part.FileName == "" {
                message.Body = append(message.Body, newBody(message, part.ContentType, part.Content))
            } else {
                message.Files = append(message.Files, part)
            }


            if err != nil {
                break
            }
        }
    } else {
        message.Body = []*Body{newBody(message, message.XHeader.Get("Content-Type"), msg.Body)}
    }
    }

    return message
}

func decodeHeaders(origin []string) []string {
    dst := make([]string, len(origin))

    for i, header := range origin {
        dst[i] = decodeHeader(header)
    }

    return dst
}

func decodeHeader(origin string) string {
    header := encodedHeaderRegexp.ReplaceAllStringFunc(origin, func(src string) string {
        if dec := encodedHeaderRegexp.FindStringSubmatch(src); len(dec) == 4 {
            return decode(dec[3], dec[1], dec[2])
        } else {
            return src
        }
    })

    return header
}

func decode(s string, charset string, encoder string) string {
    var r io.Reader
    r = strings.NewReader(s)

    dec := newDecoder(r, charset, encoder)

    dst := bytes.NewBuffer([]byte{})
    io.Copy(dst, dec)

    return dst.String()
}

func newDecoder(r io.Reader, charset string, encoder string) io.Reader {
    switch strings.ToUpper(encoder) {
    case "B", "BASE64":
        r = base64.NewDecoder(base64.StdEncoding, r)
    case "Q", "QUOTED-PRINTABLE":
        fallthrough
    default:
        r = quotedprintable.NewReader(r)
    }
    return r
}

func base64Encode(body io.Reader) string {
    w := &bytes.Buffer{}
    if _, err := io.Copy(w, body); err != nil {
        panic(err)
    }
    return base64.StdEncoding.EncodeToString(w.Bytes())
}