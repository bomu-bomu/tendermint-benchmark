# Tendermint Benchmark

## Install dependency
    ```sh
    cd $GOPATH/src/github.com/oatsaysai/tendermint-benchmark/abci
    
    dep ensure
    ```

## How to run
1.  Run ABCI server

    ```sh
    cd $GOPATH/src/github.com/oatsaysai/tendermint-benchmark
    
    go run abci/server.go
    ```
    
2.  Run tendermint

    ```sh
    tendermint init

    tendermint unsafe_reset_all && tendermint node --consensus.create_empty_blocks=false
    ```

3.  Run API

    ```sh
    cd $GOPATH/src/github.com/oatsaysai/tendermint-benchmark
    
    go run RESTAPI/main.go
    ```
    
### How to send request to Quorum API
   use [loadtest project](https://github.com/the-hulk-id/loadtest)
 ```sh
 # clone from https://github.com/the-hulk-id/loadtest
 # Example how to use loadtest
 node loadtest.js -d 10 -m -100 -s 0 -a 10.0.1.13,8181,10.0.1.14,8181,10.0.1.15,8181,10.0.1.16,8181,10.0.1.17,8181,10.0.1.18,8181
 ```

## Add new validator (For testing)
get PubKey from pub_key.data in priv_validator.json 
```sh
curl -s 'localhost:46657/broadcast_tx_commit?tx="val:PubKey"'
```
