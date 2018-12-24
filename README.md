# go-m3ua

Simple M3UA protocol implementation in pure Golang.

[![CircleCI](https://circleci.com/gh/wmnsk/go-m3ua.svg?style=svg)](https://circleci.com/gh/wmnsk/go-m3ua)
[![GoDoc](https://godoc.org/github.com/wmnsk/go-m3ua?status.svg)](https://godoc.org/github.com/wmnsk/go-m3ua)
[![Go Report Card](https://goreportcard.com/badge/github.com/wmnsk/go-m3ua)](https://goreportcard.com/report/github.com/wmnsk/go-m3ua)
[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/wmnsk/go-m3ua/blob/master/LICENSE)

## Quickstart

### Installation

The following packages should be installed before starting.

```shell-session
go get -u github.com/pkg/errors
go get -u github.com/ishidawataru/sctp
go get -u github.com/google/go-cmp/cmp
go get -u github.com/pascaldekloe/goe/verify
go get -u github.com/wmnsk/go-m3ua
```

_*Non-Linux machine is NOT supported, as this package relies much on [`github.com/ishidawataru/sctp`](https://github.com/ishidawataru/sctp)._

### Running examples as it is

Working examples are available in [examples directory](./examples/).

```shell-session
cd go-m3ua/examples/client
go run m3ua-client.go --addr <dst-IP:dst-port>
```

### Write your own code

The API design is kept as similar as possible to other protocols in standard `net` package. To establish M3UA connection as client/server, you can use `Dial()` and `Listen()`/`Accept()` without caring about the underlying SCTP association, as go-m3ua handles it together with M3UA ASPSM & ASPTM procedures.

Here is an example to develop your own M3UA client using go-m3ua.

First, you need to create `*Config` used to setup/maintain M3UA connection.

```go
config := m3ua.NewClientConfig(
    &m3ua.HeartbeatInfo{
        Enabled:  true,
        Interval: time.Duration(3 * time.Second),
        Timer:    time.Duration(10 * time.Second),
    },
    0x11111111, // OriginatingPointCode
    0x22222222, // DestinationPointCode
    1,          // AspIdentifier
    params.TrafficModeLoadshare, // TrafficModeType
    0,                     // NetworkAppearance
    0,                     // CorrelationID
    []uint32{1, 2},        // RoutingContexts
    params.ServiceIndSCCP, // ServiceIndicator
    0, // NetworkIndicator
    0, // MessagePriority
    1, // SignalingLinkSelection
)
// set nil on unnecessary paramters.
config.CorrelationID = nil
```

Then, prepare network addresses and context and try to connect with `Dial()`.

```go
// setup SCTP peer on the specified IPs and Port.
raddr, err := sctp.ResolveSCTPAddr("sctp", SERVER_IPS)
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()
ctx, cancel := context.WithCancel(ctx)
defer cancel()

conn, err := m3ua.Dial(ctx, "m3ua", nil, raddr, config)
if err != nil {
    log.Fatalf("Failed to dial M3UA: %s", err)
}
defer conn.Close()
```

Now you can `Read()` / `Write()` data from/to the remote endpoint.

```go
if _, err := conn.Write(d); err != nil {
    log.Fatalf("Failed to write M3UA data: %s", err)
}
log.Printf("Successfully sent M3UA data: %x", d)

buf := make([]byte, 1500)
n, err := conn.Read(buf)
if err != nil {
    log.Fatal(err)
}

log.Printf("Successfully read M3UA data: %x", buf[:n])
```

See [example/server directory](./examples/server) for server example.

## Supported Features

### Messages

| Class    | Message                                         | Supported | Notes                                                          |
|----------|-------------------------------------------------|-----------|----------------------------------------------------------------|
| Transfer | Payload Data Message (DATA)                     | Yes       | [RFC4666#3.3](https://tools.ietf.org/html/rfc4666#section-3.3) |
| SSNM     | Destination Unavailable (DUNA)                  |           | [RFC4666#3.4](https://tools.ietf.org/html/rfc4666#section-3.4) |
|          | Destination Available (DAVA)                    |           |                                                                |
|          | Destination State Audit (DAUD)                  |           |                                                                |
|          | Signalling Congestion (SCON)                    |           |                                                                |
|          | Destination User Part Unavailable (DUPU)        |           |                                                                |
|          | Destination Restricted (DRST)                   |           |                                                                |
| ASPSM    | ASP Up                                          | Yes       | [RFC4666#3.5](https://tools.ietf.org/html/rfc4666#section-3.5) |
|          | ASP Up Acknowledgement (ASP Up Ack)             | Yes       |                                                                |
|          | ASP Down                                        | Yes       |                                                                |
|          | ASP Down Acknowledgement (ASP Down Ack)         | Yes       |                                                                |
|          | Heartbeat (BEAT)                                | Yes       |                                                                |
|          | Heartbeat Acknowledgement (BEAT Ack)            | Yes       |                                                                |
| RKM      | Registration Request (REG REQ)                  |           | [RFC4666#3.6](https://tools.ietf.org/html/rfc4666#section-3.6) |
|          | Registration Response (REG RSP)                 |           |                                                                |
|          | Deregistration Request (DEREG REQ)              |           |                                                                |
|          | Deregistration Response (DEREG RSP)             |           |                                                                |
| ASPTM    | ASP Active                                      | Yes       | [RFC4666#3.7](https://tools.ietf.org/html/rfc4666#section-3.7) |
|          | ASP Active Acknowledgement (ASP Active Ack)     | Yes       |                                                                |
|          | ASP Inactive                                    | Yes       |                                                                |
|          | ASP Inactive Acknowledgement (ASP Inactive Ack) | Yes       |                                                                |
| MGMT     | Error                                           | Yes       | [RFC4666#3.8](https://tools.ietf.org/html/rfc4666#section-3.8) |
|          | Notify                                          | Yes       |                                                                |

### Parameters

| Type          | Parameters                   | Supported | Notes |
|---------------|------------------------------|-----------|-------|
| Common        | INFO String                  | Yes       |       |
|               | Routing Context              | Yes       |       |
|               | Diagnostic Information       | Yes       |       |
|               | Heartbeat Data               | Yes       |       |
|               | Traffic Mode Type            | Yes       |       |
|               | Error Code                   | Yes       |       |
|               | Status                       | Yes       |       |
|               | ASP Identifier               | Yes       |       |
| M3UA-specific | Network Appearance           | Yes       |       |
|               | User/Cause                   | Yes       |       |
|               | Congestion Indications       | Yes       |       |
|               | Concerned Destination        | Yes       |       |
|               | Routing Key                  | Yes       |       |
|               | Registration Result          | Yes       |       |
|               | Deregistration Result        | Yes       |       |
|               | Local Routing Key Identifier | Yes       |       |
|               | Destination Point Code       | Yes       |       |
|               | Service Indicators           | Yes       |       |
|               | Originating Point Code List  | Yes       |       |
|               | Protocol Data                | Yes       |       |
|               | Registration Status          | Yes       |       |
|               | Deregistration Status        | Yes       |       |

## Disclaimer

This is still experimental project.
In some part, the behavior is not fully compliant with RFC, and some of the features are not even implemented yet.

Also note that some exported APIs may be changed without notice before first release.

## Author

Yoshiyuki Kurauchi ([GitHub](https://github.com/wmnsk/) / [Twitter](https://twitter.com/wmnskdmms))

## LICENSE

[MIT](https://github.com/wmnsk/go-m3ua/blob/master/LICENSE)
