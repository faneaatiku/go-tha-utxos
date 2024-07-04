## Usage

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
