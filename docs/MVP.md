# MVP (1–2 weeks)

## Goal

Ship a usable **CLI** that a stranger can run in &lt; 5 minutes on an API they own and get a clear pass/fail on a quota-style invariant.

## In scope

1. **Config file** (YAML or JSON), e.g.:
   - `method`, `url`, `headers`, `body`
   - `concurrency` (N)
   - `attempts` (repeat bursts)
   - `success_rule` (e.g. max K responses with status 200)
   - optional `follow_up` GET to read remaining quota
2. **Burst runner** — fire N requests with tight synchronization (goroutines + wait group / barrier)  
3. **Report** — stdout + optional JSON:
   - counts by status  
   - whether success_rule held  
   - sample bodies (truncated) on failure  
4. **Safety banners** — “only against APIs you own / staging”  
5. **README** — quickstart, GIF/asciinema, ethics  
6. **License** — MIT or Apache-2.0  
7. **Example** — against a tiny local demo server (included) that *intentionally* has a racey counter

## Out of scope (MVP)

- Cloud SaaS UI  
- Auto-discovery of OpenAPI  
- HTTP/2 single-packet attacks (nice later; PortSwigger-style)  
- GUI  
- Paid features  
- WhatsApp / other product integrations  

## Demo server (recommended)

`examples/racey-counter`: in-memory `remaining` with naive check-then-act so MVP always has a failing case to show.

## CLI sketch (non-final)

```bash
quota-race run -c examples/wallet-quota.yaml
quota-race run -c examples/wallet-quota.yaml --json
```

## Success criteria for “MVP done”

- [x] Fresh machine: install → run example → see FAIL on racey demo
- [x] Same config against fixed (mutex/atomic) demo → PASS
- [x] README explains difference vs k6 in 5 lines
- [x] Ethics section visible
- [x] Tag `v0.1.0` on GitHub

## Suggested implementation order

1. Demo servers (racey + fixed)  
2. Burst client + report  
3. Config schema  
4. README + GIF  
5. GitHub Action smoke test on demo
