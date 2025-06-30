# blockshare
                           ┌──────────────────────┐
                           │   MCR Blockchain     │
                           │  (Golang-based PoS)  │
                           └─────────┬────────────┘
                                     │
       ┌─────────────────────────────┴─────────────────────────────┐
       │                                                           │
┌──────▼───────┐       ┌──────────────▼─────────────┐      ┌───────▼────────┐
│ Ride Module  │       │  Staking & Validator Logic │      │ Identity Module│
│ (Transactions│       │  (PoS + Physical Work)     │      │ (KYC, License, │
│ + Payments)  │       └──────────────┬─────────────┘      │ Insurance, etc)│
└──────┬───────┘                      │                    └──────┬─────────┘
       │                              │                           │
┌──────▼────────┐     ┌───────────────▼────────────┐   ┌─────────▼─────────────┐
│ Smart Contract│     │ Validator Reputation Engine │   │ External Verification │
│ (in Golang)   │     │ (Ride Count, Reviews, etc)  │   │ APIs or Admin Tools   │
└───────────────┘     └───────────────┬────────────┘   └───────────────────────┘
                                     │
                              ┌──────▼───────┐
                              │ Block Commit │
                              │ (via PoS +   │
                              │ physical rep)│
                              └──────────────┘

| Concept                                           | Component                                         | Notes                                    |
| ------------------------------------------------- | ------------------------------------------------- | ---------------------------------------- |
| Rides generate transaction data                   | `Ride Module`                                     | Ride start/stop, distance, fare          |
| Payment flows on-chain                            | `Ride Module + Smart Contract`                    | Token-based (or fiat-linked)             |
| Drivers/riders earn validator status via activity | `Staking & Validator Logic` + `Reputation Engine` | Based on number of rides, quality, trust |
| Validator checks: License, Insurance, etc.        | `Identity Module` + `External APIs`               | Stored off-chain w/ on-chain references  |
| Consensus                                         | `PoS + Ride Reputation = Validator`               | Validators stake both tokens & work      |

🔄 Current Capabilities:
✅ Genesis Validator Onboarding
BecomeValidator() lets the genesis driver become a validator with no min stake (bootstrapping phase).

✅ Ride Lifecycle Flow
Submit Ride: SubmitRideTx() from the ride module to blockshare

Verify Pickup: SubmitPickupProof() using passenger-provided code

Verify Dropoff: SubmitDropoff() with location + timestamp

Approval & Finalization: ApproveRideTx() to commit to ledger