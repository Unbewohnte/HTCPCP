# HTCPCP (RFC 2324)
## Server and client (incomplete) implementation


### Info

This is an incomplete HTCPCP [RFC 2324](https://datatracker.ietf.org/doc/html/rfc2324) implementation for the server and the client.

PROPFIND, GET, BREW and WHEN requests are handled. Additions are not supported.

### Server

On BREW and WHEN requests server *executes specified commands* from configuration file. This way it is possible for the server to actually launch some work via a script when requests come.

If configuration file is not present - the next launch will create it.

Default configuration file structure:

```json
{
 "commands": {
  "BrewCommand": "./brew.sh",
  "StopPouringCommand": "./stopPouring.sh"
 },
 "coffee-type": "Latte",
 "brew-time-sec": 10,
 "max-pour-time-sec": 5
}
```

`brew-time-sec` is the approximate time it takes to brew coffee in seconds. After this amount is passed, status of the pot changes to `Pouring` and will stay it until `max-pour-time-sec` seconds are passed or an incoming request with `WHEN` method is presented.

By default port 80 is used, but can be changed with a `-port` flag:

`HTCPCP-server -port 8000`

Server will listen and handle incoming requests.

### Client

Client is used to interact with the server and has such syntax:

```
HTCPCP-client (-version) [ADDR] [COMMAND]
  -version
        Print version information and exit
  ADDR string
    	Address of an HTCPCP server with port (ie: http://111.11.111.1:80 or http://coffeeserver:80)
  COMMAND string
    	Command to send (ie: GET, WHEN, BREW, PROPFIND)
```

Examples:

- `HTCPCP-client -help` - prints above syntax message 
- `HTCPCP-client -version` - prints version information
- `HTCPCP-client http://localhost:8000 propfind` - sends a PROPFIND request to localhost:8000
- `HTCPCP-client http://localhost:8000 brew` - sends a BREW request to localhost:8000
- `HTCPCP-client http://localhost:8000 get` - sends a get request to localhost:8000


### License

MIT License