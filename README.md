# Homevision Client

### Author: Wenceslao Marquez Burgos (Lao)

## Requirements:
- [Docker](https://docs.docker.com/engine/install/)
- [Task](https://taskfile.dev/installation/)


## HOW TO build

After installing the two dependencies all you need to do is:

1. Clone the repo `git clone https://github.com/wmb1207/homevision-client.git`
2. `cd homevision-client`
3. `task build-linux-amd64` if you are using linux or `task build-macos-arm64` If you are using a mac with M1
4. `./bin/client`


## How to run the test suite
`test-docker` This will never call the real api
`test-real-docker` This will call the real api

