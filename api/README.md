> Cache in memory with gRPC

## Contents

- [About](#about)
- [Running Server](#running-server)
- [Testing](#testing)
- [Development](#development)
- [Example](#example)

## About

Cache in memory with gRPC generate API:
- Add key/value
- Get value using a key
- Increment value using a key
- Calls can be made using goroutines on the client side.(concurrency safe)

### Running Server

You can run application locally

```bash
start server `./mockredis-server` or `make server`
`./mockredis-server -addr=":12345"` to run server on port `12345`
```

## Testing

Run `make test` to run the tests. 

```bash
go test internal/server/* -v -cover -race
=== RUN   TestSet
--- PASS: TestSet (0.02s)
=== RUN   TestDump
--- PASS: TestDump (0.00s)
=== RUN   TestIncrWithValidKey
--- PASS: TestIncrWithValidKey (0.00s)
=== RUN   TestIncrWithInvalidKey
--- PASS: TestIncrWithInvalidKey (0.00s)
PASS
coverage: 74.6% of statements
ok      command-line-arguments  0.361s  coverage: 74.6% of statements
```

## Development

```bash
# Server
api $ ./mockredis-server --help
Usage of ./mockredis-server:
  -addr string
        Address on which you want to run server (default ":12345")
  -cln int
        Cleanup interval duration of expired cache is 5 min (default 5)
  -exp int
        Default expiration duration of cache is 10 min (default 10)

# Client
api $ ./mockredis-cli --help
Usage of ./mockredis-cli:
  -addr string
        Server connection address (default "127.0.0.1:12345")
  -k string
        Key
  -o string
        Operation type: Set|Dump|Incr
  -v string
        Value
```

## Example

```bash
# Set
api $ ./mockredis-cli -o Set -k company -v atato
INFO[0000] App mockredis-cli start...                   
INFO[0000] Server connection address: 127.0.0.1:12345   
INFO[0000] Operation in progress 'OpTypeSet'            
INFO[0000] Call method SET                              
INFO[0000] Response from server: key:"company"  value:"atato"  expiration:"1m"

# Dump
api $ ./mockredis-cli -o Dump -k company 
INFO[0000] App mockredis-cli start...                   
INFO[0000] Server connection address: 127.0.0.1:12345   
INFO[0000] Operation in progress 'OpTypeDump'           
INFO[0000] Call method DUMP                             
INFO[0000] Response from server: key:"company"  value:"atato"  expiration:"2022-06-19 21:59:23.952336 +0300 MSK"

# Incr witn key and value is not integer
api $ ./mockredis-cli -o Incr -k company 
INFO[0000] App mockredis-cli start...                   
INFO[0000] Server connection address: 127.0.0.1:12345   
INFO[0000] Operation in progress 'OpTypeIncr'           
INFO[0000] Call method INCR                             
ERRO[0000] Startup error: rpc error: code = Unknown desc = value is not an integer 

# Incr with key and key is not exist
api $ ./mockredis-cli -o Incr -k var 
INFO[0000] App mockredis-cli start...                   
INFO[0000] Server connection address: 127.0.0.1:12345   
INFO[0000] Operation in progress 'OpTypeIncr'           
INFO[0000] Call method INCR                             
INFO[0000] Response from server: success:true
```