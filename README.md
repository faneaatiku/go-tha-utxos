## Usage

### Auto-pilot  
You can create UTXOs by just running this software with a simple config file:
```yaml
commands:
  daemon_cli: './tha-cli'
logging:
  level: 'info'
rpc:
    host: "{{RPC_HOST:RPC_PORT}}"
    user: "{{RPC_USER}}"
    password: "{{RPC_PASSWORD}}"
    wallet_name: "" #optional - exclude it if you use only 1 loaded wallet
auto_runner:
  addresses_file: "" #optional - if you want to use an existing file. otherwise it will create one for you with addresses fetched from the daemon
  addresses_count: 500 #optional - how many addresses to use when creating UTXOs
  utxos_interval: 2 #optional, interval to run the utxo create in minutes
```
Place this config in the same folder as the binary downloaded from this repo.

Generate new addresses and optionally save them to a file
```
./go-tha-utxos addresses generate --count 50 --file addresses.json 
```

Collect addresses from listaddressgroupings and create new one if needed  
Attention: listaddressgroupings returns addresses the addresses that were used at least once
```
./go-tha-utxos addresses collect --count 50 --file collected_addreses.json
```
If not enough addresses were found the above command creates some addresses to reach the requested number (--count flag)

Create UTXOs
- command that reads unspent outputs bigger than 0.1 and sends them to the provided addresses
- The remaining is sent back to one of the addresses from the inputs
- the fee represents the total fee paid for the entire transaction
```
./go-tha-utxos generate --file collected_addresses.json --fee 0.01
```

## EXPERIMENTAL SOFTWARE! USE AT YOUR OWN RISK
