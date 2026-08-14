# quota-race

Open-source **quota / race checker** for HTTP APIs **you own**.

It fires near-simultaneous requests and checks a **business invariant** (at most N successful grants). It is not k6. It is not for third-party production.

> Use only against APIs you own, loopback, or staging you have written permission to test.

Related lab note: [Green tests, forty grants](https://riopermana.com/writing/green-tests-forty-grants/).

## Install

```bash
go install github.com/rioaguspermana/quota-race/cmd/quota-race@latest
```

From this tree:

```bash
go test ./...
go build -o quota-race.exe ./cmd/quota-race
```

## Quickstart (demo)

Terminal 1 — intentional TOCTOU counter:

```bash
go run ./examples/counter -addr 127.0.0.1:8090 -racey -limit 10
```

Terminal 2:

```bash
go run ./cmd/quota-race run -c examples/racey.yaml
```

Expect **FAIL**: 40 concurrent takes vs limit 10 on a naive check-then-act handler.

Fixed (mutex) counter on another port:

```bash
go run ./examples/counter -addr 127.0.0.1:8091 -limit 10
go run ./cmd/quota-race run -c examples/fixed.yaml
```

Expect **PASS**.

Exit code `1` means the invariant failed (useful in CI). Exit `2` is usage/config/safety.

## Config

See `examples/racey.yaml`. Important fields:

- `request.url` / `method` / optional `headers` and `body`
- `concurrency` (hard cap 256)
- `invariant.ok_status` + `invariant.max_ok`
- `i_own_this_api: true` **required** for non-loopback URLs (or pass `--i-own-this-api`)

Loopback (`127.0.0.1`, `localhost`) does not need the flag. The ethics banner still prints.

## vs k6

k6 measures load (RPS, latency). This tool asks: **did 40 overlapping grants break a limit of 10?** Sequential tests stay green; production overlap does not.

## Ethics

See [docs/ETHICS.md](docs/ETHICS.md). This project does not authorize testing anyone else’s systems.

## License

MIT
