# SHH
Tiny CLI helper for SSH: save hosts in YAML, connect by name, and use your system ssh.
## Install
### via script

    curl -fsSL https://raw.githubusercontent.com/unholyFigaro/shh/main/install.sh | bash

### via Go

    go install github.com/unholyFigaro/shh@latest
## Usage

    # add
    shh add dev-stand --host 127.0.0.1 -p 2222 -u ubuntu
    
    # list
    shh list
    
    # connect
    shh dev-stand
    
## Config
Default: `~/.config/shh/hosts.yaml` (or set `SHH_CONFIG`).

    version: "1.0"
    hosts:
      dev-stand:
        host: 127.0.0.1
        port: 2222
        user: ubuntu
## License
0BSD — do anything; no warranty. See `LICENSE`.

