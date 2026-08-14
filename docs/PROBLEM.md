# Problem

## The bug pattern

Many APIs enforce limits with **check-then-act**:

1. `SELECT remaining` / read counter  
2. If `remaining >= 1` → allow  
3. Later `UPDATE used = used + 1`

Under **sequential** tests (local, Postman, happy path), this looks fine.  
Under **concurrent** requests, two handlers can both pass the check before either writes → **quota breach**, negative stock, double trial redeem, etc.

Classic name: **TOCTOU** (Time-of-Check to Time-of-Use).  
Related: race condition on shared mutable state / business invariant.

## Why existing tests miss it

| Approach | What it catches | What it often misses |
|----------|-----------------|----------------------|
| Unit / sequential E2E | Logic for one request | Overlap windows |
| Postman collections | Status / schema | Concurrent correctness |
| Load tools (k6, etc.) | Latency, RPS, errors under load | “Did invariant break?” unless you script it |

## Who feels the pain

- Backend / fullstack engineers shipping wallet, seats, coupons, free-tier quotas  
- Teams that “passed QA” then burned credits in production  
- Security-minded engineers (race classes show up in bounty writeups)

## Job to be done

> “Before release, prove that my limit still holds when several clients hit the same endpoint at once.”

## Non-goals (problem framing)

- Replacing APM / observability  
- Guaranteeing security of third-party APIs  
- Teaching distributed consensus from scratch
