# Go Blockchain

## Panoramica

Progetto per l'esame pratico di Sicurezza dell'Informazione M.
Questo progetto educativo implementa una semplice blockchain in Go, progettata per comprendere i concetti fondamentali delle criptovalute e della tecnologia blockchain. Il sistema include proof of work, transazioni e una CLI interattiva per l'uso pratico.

## Architettura del Sistema

### Componenti Principali

```
go-blockchain/
├── main.go                 # Entry point dell'applicazione
├── cli/
│   └── cli.go             # Command Line Interface
├── blockchain/
│   ├── block.go           # Struttura dei blocchi
│   ├── blockchain.go      # Logica della blockchain
│   ├── proof.go           # Proof of Work
│   └── transaction.go     # Sistema delle transazioni
├── wallet/                # Modulo wallet (futuro)
└── tmp/
    └── blocks/           # Database BadgerDB
```

## Scelte Progettuali

### 1. **Linguaggio di Programmazione - Go**
**Motivazione:**
- Sintassi semplice e chiara
- Linguaggio moderno e performante
- Performance elevate per operazioni crittografiche
- Concurrency nativa per future implementazioni di rete

### 2. **Database - BadgerDB**
**Motivazione:**
- Database key-value embedded, no setup esterno richiesto
- Performance ottimali per letture/scritture sequenziali

### 3. **Proof of Work (PoW)**
**Implementazione:**
```go
type ProofOfWork struct {
    Block  *Block
    Target *big.Int
}
```

### 4. **UTXO Model (Unspent Transaction Outputs)**
**Motivazione:**
- Modello più sicuro rispetto agli account-based
- Prevenzione naturale del double-spending
- Tracciabilità completa delle transazioni
- Parallelizzazione delle validazioni

## Implementazioni Chiave

### Struttura dei Blocchi

```go
type Block struct {
    Timestamp    int64
    Hash         []byte
    PrevHash     []byte
    Target       []byte
    Nonce        int
    Transactions []*Transaction
}
```

### Sistema delle Transazioni

#### Transazioni Coinbase
```go
func CoinbaseTx(to, data string) *Transaction {
    txin := TxInput{[]byte{}, -1, data}
    txout := TxOutput{100, to}
    // ...
}
```
**Caratteristiche:**
- Input vuoto (ID: [], Out: -1) per identificare transazioni di mining
- Reward fisso di 100 coins per semplicità
- Utilizzate per introdurre nuovi coins nel sistema

#### Transazioni Standard
```go
func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction
```
**Features:**
- Validazione automatica dei fondi disponibili
- Gestione del "resto" (change output)
- Riferimenti agli UTXO precedenti

### Proof of Work

#### Algoritmo di Mining
```go
func (pow *ProofOfWork) Run() (int, []byte) {
    var intHash big.Int
    var hash [32]byte
    nonce := 0
    
    for nonce < maxNonce {
        data := pow.InitData(nonce)
        hash = sha256.Sum256(data)
        intHash.SetBytes(hash[:])
        
        if intHash.Cmp(pow.Target) == -1 {
            break
        }
        nonce++
    }
    return nonce, hash[:]
}
```

**Caratteristiche:**
- Difficoltà configurabile tramite target
- SHA-256 come funzione hash
- Incremento del nonce per trovare hash validi


## Comandi Disponibili
| Comando | Descrizione | Esempio |
|---------|-------------|---------|
| `createblockchain` | Crea nuova blockchain | `createblockchain -address alice` |
| `mine` | Mina nuovo blocco | `mine -address alice` |
| `send` | Invia transazione | `send -from alice -to bob -amount 50` |
| `getbalance` | Mostra bilancio | `getbalance -address alice` |
| `printchain` | Stampa tutti i blocchi | `printchain` |
| `interactive` | Modalità interattiva | `interactive` |

## Testing e Utilizzo

### Scenario di Test Completo
```bash
# 1. Crea blockchain
./go-blockchain.exe createblockchain -address genesis

# 2. Mina coins per alice
./go-blockchain.exe mine -address alice

# 3. Verifica bilancio
./go-blockchain.exe getbalance -address alice
# Output: Balance of alice: 100

# 4. Trasferisci fondi
./go-blockchain.exe send -from alice -to bob -amount 30

# 5. Verifica bilanci finali
./go-blockchain.exe getbalance -address alice  # 70
./go-blockchain.exe getbalance -address bob    # 30

# 6. Visualizza blockchain
./go-blockchain.exe printchain
```
