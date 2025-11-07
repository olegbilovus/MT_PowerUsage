# MT_PowerUsage

A Go application for monitoring and recording the power consumption of smart plugs,
supporting Shelly devices, and storing results in SQLite.

- **Smart Plug Integration**:  
  Supports different versions/models of Shelly smart plugs (`ShellyPlugS`, `ShellyPlugSv2`), with
  capability for easy extension.
- **Container-Ready**:  
  Supplied with a `Dockerfile` and `compose.yml` for straightforward reproducible deployments.

---

## Quick Start

### Prerequisites

- Go 1.25+ (for local builds, otherwise use Docker)
- Docker & Docker Compose (recommended for deployment)
- A Shelly plug accessible via IP

### Environment Variables

The application is configured by the following variables:

- `PLUG_IP`: IP address of the Shelly smart plug (required)
- `PLUG_TYPE`: 1 for ShellyPlugS, 2 for ShellyPlugSv2 (required)
- `FREQ`: Polling frequency in **seconds** (required, e.g., `1`)
- `SQLITE_DB`: SQLite database file name (required, e.g., `power_data.db`)
- `RESET_DB`: Set to `true` to drop and recreate DB table at startup (optional)

You can supply these variables via a `.env` file or directly into the environment.

### Running with Docker Compose

1. Create a `.env` file:
    ```
    PLUG_IP=192.168.1.10
    PLUG_TYPE=1
    SQLITE_DB=power_data.db
    ```

2. Start with Docker Compose:
    ```bash
    docker compose up
    ```

3. Data will be saved into `./dbs/<SQLITE_DB>` on your host machine.

### Local Run (Go)

```bash
export PLUG_IP=192.168.1.10
export PLUG_TYPE=1
export FREQ=1
export SQLITE_DB=power_data.db
go run ./cmd/MT_PowerUsage
```

---

## Database

Data is stored in an SQLite table `power` (schema auto-managed):

| timestamp (PK) | load (float)           |
|----------------|------------------------|
| TIMESTAMP      | Instantaneous power(W) |


---

## Extending

To add new device types, implement the `Plug` interface in `pkg/plug/plug.go`.

---

## License

MIT License. See [LICENSE](LICENCE).
