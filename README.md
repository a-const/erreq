
# bbreq

Little tool to make requests to Prysm gRPC easier.

# Usage

```
./brreq [-d={value in seconds}] [endpoint] [flags]

```
For Example:
```
./brreq -d=1 validators -id=head -v=10
./brreq block -id=head
```

Available endpoints:

//Node  
peers  
peerbyid  
syncing  
identity  
peer_count  
version  

//Beacon  
genesis  
validators  
root  
fork  
finality_checkpoints  
validator_by_id  
block_by_id  
