# brreq (beacon-chain rpc requester)

Little tool to make requests to Prysm gRPC easier

ver 0.2

---

Usage:

```
./brreq [-d] [command] [-key] ...
```

### Global keys

-d  Turns on live mode and sets delay for refresh

### Commands:

##### Node

* peers
  *-id*    get peer by id
* syncing
* identity
* peer_count
* version

##### Beacon

* genesis  
* validators  
  *-id*     state id (required)  
  *-v*      validator index  
* root
  *-id*     state id (required)  
* fork  
  *-id*     state id (required)  
* finality_checkpoints  
  *-id*     state id (required)  
* block  
  *-id*     block id (required)  

##### Analyse

* prop-count  
  *-from*     begin slot number (required)  
  *-to*       end slot number (or just "head") (required)  
  *-filename*     name for output file (default - "proposers.json" to same dir as bin location)  
