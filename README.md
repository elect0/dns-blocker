# Go DNS blocker

a fast, configurable DNS ad-blocking server written in Go.

1.  **Clone:**
    ```bash
    git clone [https://github.com/your-username/go-dns-blocker.git](https://github.com/your-username/go-dns-blocker.git)
    cd go-dns-blocker
    ```

2.  **Configure:**
    ```bash
    touch config.yaml
    ```
    *Config example*
    ```yaml
    # Configuration for the Go DNS Blocker
    # Copy this file to 'config.yaml' and edit it to your needs.
    listenAddress: ":53"
  
    upstreamServer: "https://1.1.1.1/dns-query"
  
    blocklistPath: "blocklist.txt"
  
    customRecords:
    "nas.local": "192.168.1.100"
    "dev.local": "127.0.0.1"
    
    logging:
    level: "info" # Options: "debug", "info", "warn", "error"
    ```

3.  **Create Blocklist:**
    ```bash
    touch blocklist.txt
    ```
    *(Add domains to `blocklist.txt`, one per line.)*

4.  **Build:**
    ```bash
    go build -o dns-blocker ./cmd/dns-blocker
    ```

## Run

Run the server with `sudo` to use the standard DNS port (53).

```bash
sudo ./dns-blocker
