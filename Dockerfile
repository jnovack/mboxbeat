FROM golang:1.10-alpine3.8 AS builder
RUN apk add --update git make
WORKDIR /go/src/github.com/jnovack/mboxbeat
COPY . /go/src/github.com/jnovack/mboxbeat
RUN make install_deps && make install

FROM alpine:3.8

ARG version=0.0.0-local
ARG build_date=unknown
ARG commit_hash=unknown
ARG vcs_url=unknown
ARG vcs_branch=unknown

ENTRYPOINT ["/usr/bin/mboxbeat"]

LABEL org.label-schema.vendor='Justin J. Novack' \
    org.label-schema.name='mboxbeat' \
    org.label-schema.description='inject /var/spool/mail/user messages into elasticsearch' \
    org.label-schema.usage='https://github.com/jnovack/mboxbeat/blob/master/docs/Programming-Guide.md' \
    org.label-schema.url='https://github.com/jnovack/mboxbeat' \
    org.label-schema.vcs-url=$vcs_url \
    org.label-schema.vcs-branch=$vcs_branch \
    org.label-schema.vcs-ref=$commit_hash \
    org.label-schema.version=$version \
    org.label-schema.schema-version='1.0' \
    org.label-schema.build-date=$build_date

COPY --from=builder /go/bin/mboxbeat /usr/bin/mboxbeat