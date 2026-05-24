# go-user-service — Note di studio

## Obiettivo del progetto

Progetto backend Go per consolidare i concetti fondamentali in vista di colloqui tecnici.  
Architettura classica a tre layer: **handler → service → repository**.  
Il focus non è solo "farlo funzionare" ma scrivere Go idiomatico, professionale, e saperlo spiegare.

---

## Stato attuale

### Completato

- `internal/model/model.go` — struct `User` con ID, Name, Email, CreatedAt, UpdatedAt
- `internal/dto/user_request.go` — `CreateUserRequest` e `UpdateUserRequest`
- `internal/errors/errors.go` — errori custom: `ErrUserNotFound`, `ErrInvalidInput`, `ErrEmailAlreadyExists`
- `internal/repository/user_repository.go` — interfaccia `UserRepository`
- `internal/repository/in_memory_user_repository.go` — struct `InMemoryUserRepository` + costruttore `NewInMemoryUserRepository()`

### In corso

Implementazione dei metodi di `InMemoryUserRepository`:
- [ ] `Create`
- [ ] `GetByID`
- [ ] `GetAll`
- [ ] `Update`
- [ ] `Delete`

### Da fare

- `internal/service/` — business logic
- `internal/handler/` — HTTP handlers
- `cmd/api/main.go` — entry point, wiring dei layer

---

## Concetti chiave già affrontati (da saper spiegare al colloquio)

### Zero values
In Go ogni tipo ha un valore di default. `int` → `0`, `string` → `""`, `sync.RWMutex` → unlocked e pronto.  
La **map** invece ha zero value `nil` — scriverci causa panic a runtime. Va inizializzata con `make()`.  
Regola pratica: inizializza solo ciò che non può stare al suo zero value utile.

### map[int]model.User vs map[int]*model.User
Le map in Go restituiscono **copie** del valore, non riferimenti.  
Modificare il valore estratto non cambia la map — devi riassegnare esplicitamente.  
Questo protegge i dati interni da mutazioni esterne involontarie → vantaggio per un repository.

### sync.RWMutex vs sync.Mutex
`sync.Mutex` → una goroutine alla volta, chiunque aspetta.  
`sync.RWMutex` → più lettori contemporanei (`RLock`), un solo scrittore (`Lock`).  
Nel repository: `GetByID` e `GetAll` usano `RLock`, `Create`/`Update`/`Delete` usano `Lock`.  
**Mai copiare un mutex** — `go vet` lo segnala come errore. Per questo il costruttore restituisce `*InMemoryUserRepository`.

### Costruttore idiomatico Go
Go non ha costruttori. La convenzione è `NewXxx() *Xxx`.  
Serve per garantire che l'istanza nasca in uno stato valido (es. map inizializzata).  
Restituisce sempre un puntatore per: (1) non copiare il mutex, (2) condividere lo stesso stato.

### Exported vs Unexported
Maiuscola = esportato (visibile fuori dal package).  
Minuscola = unexported (dettaglio implementativo, nascosto).  
I campi interni di una struct che non devono essere accessibili dall'esterno vanno in minuscolo.

### Repository pattern
Il repository è responsabile della persistenza e di tutto ciò che riguarda i dati:
- Generazione degli ID (il chiamante non sa come funzionano)
- Impostazione di `CreatedAt` / `UpdatedAt`
- Gestione della concorrenza (mutex)

Il service non sa se i dati stanno in memoria, Postgres o Redis. Parla solo con l'interfaccia.

---

## Prossimo step concreto

Scrivere il metodo `Create` su `InMemoryUserRepository`:

1. `r.mu.Lock()` + `defer r.mu.Unlock()`
2. Incrementare `r.counterId`, assegnarlo come `user.ID`
3. `user.CreatedAt = time.Now()` e `user.UpdatedAt = time.Now()`
4. `r.users[user.ID] = user`
5. `return nil`

Poi procedere con `GetByID` (usa `RLock`), `GetAll`, `Update`, `Delete`.

---

## Domande da saper rispondere a un colloquio

- Perché usi `sync.RWMutex` invece di `sync.Mutex`?
- Cosa succede se copi una struct che contiene un mutex?
- Perché il costruttore restituisce un puntatore?
- Cosa ritorna una map se la chiave non esiste?
- Differenza tra `make` e `new` in Go?
- Come implementeresti lo stesso repository con un database reale?
- Perché separare interface e implementazione in package diversi?
