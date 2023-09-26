# erreq (eth rpc requester)

Little tool to make requests to Prysm and Geth gRPC easier

ver 0.3

---

Usage:

```
./erreq [-d] [command] [-key] ...
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
  *-s*     state id (required)   
  *-v*      validator index
* root   
  *-s*     state id (required)   
* fork   
  *-s*     state id (required)  
* finality_checkpoints  
  *-s*     state id (required)  
* block   
  *-id*     block id (required)

##### Analyse

* prop-count  
  *-f*     begin slot number (required)  
  *-t*       end slot number (or just "head") (required)    
  *-filename*     name for output file (default - "proposers.json" to same dir as bin location)