# 🧱 Blockshare — Proof-of-Stake Blockchain for Rideshare

**Blockshare** is the decentralized transaction ledger powering [Music City Rideshare](https://www.musiccityrideshare.com).  
Built in Go, it introduces real-world proof-of-stake via physical work — like verified ride activity — and driver participation.

---

## 🚗⛓ What It Does

- Validates & commits ride transactions (`RideTx`)
- Enforces pickup/drop-off proof via secure codes
- Allows trusted drivers to stake tokens & become validators
- Tracks validator approvals and builds a ledger of completed rides

---

## ✨ Why This Exists

Current rideshare platforms take huge cuts.  
Blockshare empowers **drivers**, not platforms, to verify and record rides on-chain.

This is part of a broader project:  
🛠️ [musiccityrideshare.com](https://www.musiccityrideshare.com)

---

## 🧪 Features

- 🪙 Staking and validator registration
- 🔐 Pickup proof with confirmation code
- 📦 JSON-based ride ledger (local block storage)
- ⛓️ Quorum-based transaction approvals
- ⚡ Lightweight Go module (no server needed)

---

## 🔧 Usage

Import into any Go app:

```go
import "github.com/x-MrPhillips-x/blockshare/blockchain"

rc, _ := blockchain.NewRideChain("path/to/token_ledger.json")

tx := blockchain.RideTx{...}
txID, err := rc.SubmitRideTx(tx)
```

To run tests:

```bash
go test ./...
```
📜 License
MIT