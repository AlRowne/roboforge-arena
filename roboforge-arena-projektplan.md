# RoboForge Arena – Projektplan

## 1. Projektziel

**RoboForge Arena** ist ein portfoliofähiges Go-Backend, in dem Nutzer modulare Kampfroboter konfigurieren und asynchron gegeneinander antreten lassen.

Der Server:

- verwaltet Nutzer und Roboter,
- nimmt Match-Aufträge über eine REST API an,
- verarbeitet Matches in einem Hintergrund-Worker,
- simuliert Kämpfe reproduzierbar,
- speichert Ergebnis und rundenweises Event-Log,
- berechnet eine einfache Rangliste.

Eine kleine Go-CLI dient als Client und macht das Projekt ohne großes Web-Frontend vorführbar.

### Hauptlernziel

Der Schwerpunkt liegt nicht auf neuer Go-Syntax, sondern auf dem Übergang von geführten Projekten zu selbstständiger Entwicklung:

> Anforderungen verstehen → Entwurf begründen → implementieren → testen → betreiben

---

## 2. Leitidee und Abhängigkeiten

Wir beginnen nicht mit Auth und gewöhnlichem CRUD. Zuerst prüfen wir die zwei besonderen und riskanten Teile des Projekts:

1. eine deterministische Kampfsimulation,
2. das sichere Beanspruchen eines Jobs durch genau einen Worker.

Danach entsteht ein durchgängiger Ablauf:

> CLI → API → PostgreSQL → Worker → Simulation → Ergebnis → CLI

![[viz-roboforge-projektplan-1788835993252.png|700]]

### Drei Grundlagen

1. **Feste Eingaben und ein fester Seed müssen denselben Kampf erzeugen.**  
   Versteckte Zufalls- oder Zeitquellen würden reproduzierbare Tests verhindern.

2. **Zusammengehörige Datenbankänderungen bilden eine fachliche Zustandsänderung.**  
   Ergebnis, Event-Log und finaler Match-Status dürfen keinen halbfertigen Zustand hinterlassen.

3. **Prozessspeicher ist nicht zwischen Serverinstanzen geteilt.**  
   Ein Go-Mutex schützt Goroutines in einem Prozess, aber nicht mehrere Prozesse oder Container. Prozessübergreifende Job-Koordination muss deshalb über die gemeinsame Datenbank erfolgen.

---

## 3. Kern-Spielablauf

1. Ein Nutzer registriert sich und meldet sich an.
2. Er baut einen Roboter aus begrenzten Komponenten.
3. Er definiert priorisierte Verhaltensregeln.
4. Er fordert einen anderen Roboter heraus.
5. Die API legt ein wartendes Match an und gibt dessen ID zurück.
6. Ein Worker beansprucht und simuliert das Match.
7. Der Nutzer fragt Status, Ergebnis und Kampfprotokoll ab.
8. Das abgeschlossene Match fließt in die Rangliste ein.

### Ziel für die ersten Kampfregeln

Die Regeln sollen klein genug bleiben, dass ihre Korrektheit testbar ist. Eine mögliche Startgröße:

- drei Chassis-Typen,
- ungefähr fünf Module,
- eine Energie-Ressource,
- wenige Aktionen wie Angriff, Verteidigung und Reparatur,
- priorisierte Strategiebedingungen,
- feste maximale Rundenzahl,
- Sieg, Niederlage oder Unentschieden.

Die endgültigen Regeln werden vor der Implementierung gemeinsam hergeleitet und anschließend eingefroren.

---

## 4. MVP-Umfang

Ein Nutzer kann über API und CLI:

- sich registrieren und anmelden,
- einen Roboter erstellen,
- eigene und herausforderbare Roboter anzeigen,
- ein Match anlegen,
- den Match-Status abfragen,
- Ergebnis und Event-Log abrufen,
- die Rangliste anzeigen.

### Vorläufige API-Fläche

Die konkreten Pfade werden beim API-Entwurf geprüft. Als Größenordnung reichen:

- Registrierung und Login,
- Roboter erstellen und auflisten,
- Match erstellen,
- Match inklusive Status abrufen,
- Match-Events abrufen,
- Rangliste abrufen.

Ein erfolgreich angenommener, aber noch nicht simulierter Match-Auftrag kann mit `202 Accepted` und einer Match-ID beantwortet werden.

---

## 5. Zentrale Invarianten

Diese Regeln müssen unabhängig von einzelnen Beispielen immer gelten:

### Simulation

- Gleiche Bot-Versionen, gleiche Regelversion und gleicher Seed erzeugen dasselbe Ergebnis und Event-Log.
- Energie und andere begrenzte Ressourcen werden nie negativ.
- Ein zerstörter Roboter führt keine weitere Aktion aus.
- Ein Match endet spätestens nach der maximalen Rundenzahl.
- Jede gespeicherte Aktion war nach dem unmittelbar vorherigen Zustand erlaubt.

### Datenhaltung

- Ein Match verwendet unveränderliche Bot-Konfigurationen.
- Spätere Bot-Änderungen verändern alte Matches nicht.
- Ein Match hat nur gültige Zustandsübergänge:
  - `pending → running`
  - `running → completed`
  - `running → failed`
- Ergebnis, Events und finaler Status werden gemeinsam gespeichert oder gemeinsam verworfen.
- Derselbe logische Match-Request erzeugt bei einem Retry nicht unbemerkt mehrere Matches.

### Worker

- Ein Match wird zu einem Zeitpunkt von höchstens einem Worker beansprucht.
- Ein abgestürzter Worker blockiert ein Match nicht dauerhaft.
- Eine erneute Verarbeitung beschädigt keine bereits gültigen Daten.

---

## 6. Zielarchitektur

### Modularer Monolith

Ein Repository und eine gemeinsam genutzte Go-Domain-Logik. Keine Microservices.

Bestandteile:

- **API-Prozess:** HTTP, Authentifizierung und Validierung
- **Worker-Prozess:** beansprucht und simuliert wartende Matches
- **CLI:** spricht ausschließlich mit der REST API
- **PostgreSQL:** persistente Daten und einfache Job Queue
- **Simulation Engine:** kennt weder HTTP noch SQL

Die genaue Paketstruktur wird nicht vorgegeben. Sie wird anhand der Abhängigkeiten entworfen. Unterstützender interner Code gehört voraussichtlich unter `internal/`; ein allgemeines `util`-Paket soll vermieden werden.

### Datenbank als Job Queue

Für den MVP ist kein Redis oder Kafka nötig. PostgreSQL ist bereits der gemeinsame, persistente Koordinator.

Ein Worker soll Matches atomar beanspruchen können, beispielsweise mithilfe einer PostgreSQL-Sperrstrategie wie `FOR UPDATE SKIP LOCKED`. Eine zeitlich begrenzte Lease ermöglicht später die Wiederaufnahme, falls ein Worker abstürzt.

Die genaue SQL-Lösung wird zunächst als kleine Risikoprobe untersucht und nicht blind aus diesem Plan übernommen.

### Event-Log statt vollständigem Event Sourcing

Kampfereignisse werden unveränderlich gespeichert, damit ein Match nachvollzogen werden kann. Daraus wird jedoch kein vollständiges Event-Sourcing-/CQRS-System gebaut.

---

## 7. Teststrategie

Tests entstehen während der Implementierung, nicht erst am Ende.

### Simulation

- tabellengesteuerte Unit-Tests,
- fester Seed für reproduzierbare Kämpfe,
- Vergleich vollständiger Event-Sequenzen für ausgewählte Fälle,
- ein Fuzz-Test für gültige und ungültige Konfigurationen,
- Prüfung allgemeiner Invarianten.

### API

- Handler-Tests mit `httptest`,
- Authentifizierungs- und Autorisierungsfälle,
- Validierungsfehler,
- korrekte HTTP-Statuscodes,
- Idempotenz eines wiederholten Match-Requests.

### PostgreSQL und Worker

- Integrationstests gegen echtes PostgreSQL,
- mehrere konkurrierende Worker,
- Transaktions-Rollback bei einem Fehler,
- Wiederaufnahme nach abgelaufener Lease,
- Migrationen auf einer leeren Datenbank.

### CI

Mindestens:

- Format-/Vet-Prüfung,
- Unit-Tests,
- `go test -race`,
- klar getrennte Integrationstests.

---

## 8. 40-Stunden-Roadmap

| Abschnitt | Zeit | Ergebnis |
|---|---:|---|
| Produktdefinition und Invarianten | 3 h | kleine, eingefrorene Kampfregeln |
| Zwei technische Risikoproben | 4 h | deterministischer Mini-Kampf und atomarer Job-Claim |
| Kampfsimulation mit Unit-Tests | 6 h | reine, reproduzierbare Engine |
| Datenmodell, Migrationen, Auth und Bots | 5 h | persistente Nutzer und Bot-Konfigurationen |
| Match-API, Worker und Transaktionen | 7 h | vollständiger asynchroner Match-Lifecycle |
| Event-Log und einfache Rangliste | 4 h | nachvollziehbare Matches und Scores |
| Go-CLI | 3 h | demonstrierbarer End-to-End-Ablauf |
| Integration, Fuzzing und Race-Tests | 4 h | geprüfte Fehler- und Parallelitätsfälle |
| Docker, Deployment, CI und README | 4 h | reproduzierbar startbares Portfolio-Projekt |
| **Gesamt** | **40 h** | |

Die Zeiten sind ein Budget, kein Vertrag. Bei Zeitdruck wird zuerst ein Feature gekürzt, nicht Datenintegrität, Tests oder Dokumentation.

---

## 9. Meilensteine und Abnahmekriterien

### M0 – Regeln stehen fest

Fertig, wenn:

- ein Kampf auf Papier eindeutig durchgeführt werden kann,
- Ressourcen, Aktionen und Abbruchbedingungen definiert sind,
- die wichtigsten Invarianten notiert sind.

### M1 – Reine Simulation

Fertig, wenn:

- kein HTTP und keine Datenbank benötigt werden,
- derselbe Seed wiederholt dasselbe Event-Log erzeugt,
- ein Kampf immer terminiert,
- Kernregeln automatisiert getestet sind.

### M2 – Persistente API

Fertig, wenn:

- Nutzer und Bots gespeichert werden,
- Authentifizierung funktioniert,
- ein Match als `pending` angelegt und abgefragt werden kann,
- eine spätere Bot-Änderung alte Matches nicht verändert.

### M3 – Zuverlässiger Worker

Fertig, wenn:

- mehrere Worker nicht gleichzeitig dasselbe Match verarbeiten,
- ein Match `pending → running → completed/failed` durchläuft,
- ein abgestürzter Worker keinen dauerhaften Lock hinterlässt,
- Ergebnis und Events transaktional gespeichert werden.

### M4 – Vorführbares Produkt

Fertig, wenn:

- der komplette Ablauf über die CLI möglich ist,
- Event-Logs lesbar ausgegeben werden,
- eine einfache Rangliste existiert.

### M5 – Portfolioqualität

Fertig, wenn:

- das Projekt mit dokumentierten Befehlen gestartet werden kann,
- Docker Compose API, Worker und PostgreSQL startet,
- Migrationen reproduzierbar laufen,
- Tests und Race Detector in CI laufen,
- README Architektur und wichtige Trade-offs erklärt,
- ein frischer Nutzer den Demo-Ablauf nachvollziehen kann.

---

## 10. Bewusst nicht im MVP

- Ausführen von Nutzercode
- Echtzeit-Multiplayer
- grafisches Web-Frontend
- Microservices
- Redis oder Kafka
- WebSockets
- OAuth und Passwort-Reset
- Kubernetes
- GraphQL oder gRPC
- vollständiges Event Sourcing/CQRS
- komplexes Matchmaking
- umfangreiches Admin-System
- ausgefeiltes Balancing

---

## 11. Stretch Goals

Erst nach vollständiger MVP-Abnahme:

1. Elo-Rating
2. Saisons
3. Bot-Winrates und Modulstatistiken
4. tägliche Turniere
5. Server-Sent Events für Match-Status
6. zusätzliche Strategiebedingungen
7. Regelversionen mit mehreren kompatiblen Rulesets

Maximal ein oder zwei Stretch Goals auswählen.

---

## 12. Zusammenarbeit mit dem Lehrer

Für jeden Abschnitt gilt:

1. Der Lehrer stellt Problem, Ziel und unverhandelbare Anforderungen vor.
2. Der Schüler entwirft eine Lösung.
3. Der Lehrer prüft den Entwurf mit Fragen und Gegenbeispielen.
4. Der Schüler implementiert selbst.
5. Fehler werden gemeinsam analysiert.
6. Hinweise werden stufenweise gegeben; es gibt keinen fertigen Code zum Abschreiben.

Der nächste Schritt ist **noch kein Code**. Zuerst werden die kleinstmöglichen Kampfregeln und ihre Invarianten entworfen.

---

## 13. Technische Referenzen

- [Go: Organizing a Go module](https://go.dev/doc/modules/layout)
- [PostgreSQL: Explicit Locking](https://www.postgresql.org/docs/current/explicit-locking.html)
- [PostgreSQL: Transaction Isolation](https://www.postgresql.org/docs/current/transaction-iso.html)
- [Go Race Detector](https://go.dev/doc/articles/race_detector)
- [Testcontainers for Go – PostgreSQL](https://golang.testcontainers.org/modules/postgres/)
- [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [The Twelve-Factor App: Config](https://12factor.net/config)
