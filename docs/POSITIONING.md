# Positioning

## Category

**API concurrency correctness** / **business-invariant race testing**  
(Not “performance testing”, not “pentest framework” as primary pitch)

## Positioning statement

For backend engineers who ship limited resources (quotas, stock, seats),  
**quota-race** is an open-source CLI that bursts concurrent requests and checks whether your limit still holds —  
unlike load testers that optimize for RPS, and unlike sequential QA that never overlaps.

## Message house

| Pillar | Line |
|--------|------|
| Problem | Local sequential green ≠ production-safe limits |
| Mechanism | Concurrent burst + invariant assertion |
| Proof | Demo racey counter fails; fixed counter passes |
| Trust | OSS, own-API only, CI-friendly |
| Tone | Factual, humble, engineer-to-engineer (same as personal brand) |

## Vs alternatives

| Tool type | Primary job | We emphasize |
|-----------|-------------|--------------|
| k6 / Locust / Artillery | Load, latency, throughput | Invariant / quota correctness |
| Postman | Manual / sequential API checks | Concurrent overlap |
| RaceGuard / racey / etc. | Similar race detection | Clear quota/limit scenarios + app-dev DX + example servers + CI story |

Be honest in docs: stand on related tools’ shoulders; differentiate with **quota/limit templates**, **simple success rules**, **dual demos**, **Indonesian+EN docs optional**.

## Product Hunt one-liners (draft)

- “Find quota races before production does.”  
- “Concurrent requests vs your business limits — pass or fail.”  
- “TOCTOU checks for HTTP quotas, open source.”

## What not to claim

- Empty market / “first ever”  
- Guaranteed detection of all races  
- Safe to run on third-party production  
- Market size USD figures without vendor disclaimer
