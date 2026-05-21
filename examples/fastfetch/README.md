# SMS Fastfetch

CLI tool that runs `fastfetch` on the host and syncs the system info to the SMS Warehouse as a datasource.

## Prerequisites

- [`fastfetch`](https://github.com/fastfetch-cli/fastfetch) installed on the machine you want to report from
- A running SMS Warehouse reachable from that machine

## Steps

### 1. Build the binary

On any machine with Go installed:

```bash
cd examples/fastfetch
GOOS=linux GOARCH=amd64 go build -o sms-fastfetch .
```

### 2. Copy to the target machine

```bash
scp sms-fastfetch user@yourserver:~/
```

### 3. Run it

```bash
./sms-fastfetch -name "my-server" -warehouse http://<warehouse-host>:8081
```

`-name` is the datasource name that will appear in the backoffice. `-warehouse` defaults to `http://localhost:8081` if omitted.

### 4. (Optional) Run on a schedule

To keep the data fresh, add a cron job on the target machine:

```bash
crontab -e
```

Add a line to run every hour:

```
0 * * * * /home/user/sms-fastfetch -name "my-server" -warehouse http://<warehouse-host>:8081
```

## Notes

- Each run replaces all existing datapoints for the datasource, so the data is always a fresh snapshot.
- If the warehouse is only reachable via a `kubectl port-forward` on another machine, start the forward with `--address 0.0.0.0` so it listens on the LAN interface:
  ```bash
  kubectl port-forward -n monitoring svc/sms-warehouse 8081:8081 --address 0.0.0.0
  ```
