# go-erc20leaderboard

A simple tool that displays the top 5 addresses on the ETH mainnet with activity metrics for the last 100 blocks, the activity metric increases if the wallet sends or receives any ERC20 token.

# Run in docker

1. create config from example
```bash
~$ cp config.example config
~$ vi config
```

2. run build
```bash
~$ make build-amd64
```

3. build docker-image
```bash
~$ make docker-build
```

4. run
```bash
~$ make docker-run
```

# License

**MIT**
