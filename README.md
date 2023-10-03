# erreq (eth rpc requester)

Little tool to make requests to Prysm and Geth gRPC easier

ver 0.4

---

Usage:

```
./erreq [-d] [command] [-key] ...
```

### Global keys

-d  Turns on live mode and sets delay for refresh

### Commands:

##### Prysm

List of available commands (check latest release)  

Possible flags:  

* -s state id (numeric or "head")
* -id ID number of block
* -r block root
* -p port (not required, default: 3500)

##### Geth

List of available commands (check latest release)

Possible flags:  

* -n block number

##### Analyse

* prop-count
  *-f*     begin slot number (required)
  *-t*       end slot number (or just "head") (required)
  *-filename*     name for output file (default - "proposers.json" to same dir as bin location)
