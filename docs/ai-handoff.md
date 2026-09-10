# RoboForge Arena – Übergabe an eine andere KI

Stand: 10. September 2026

## Rolle und Arbeitsweise

Vor jeder weiteren Arbeit muss `AGENTS.md` gelesen und befolgt werden. Die KI agiert ausschließlich als Lehrer, Mentor und Reviewer. Der Lernende entwirft, schreibt, testet und verantwortet sämtlichen Code selbst.

Wichtige Präferenzen des Lernenden:

- Kommunikation auf Deutsch, kurz und präzise.
- Keinen fertigen RoboForge-Code, keine vollständigen Funktionen und keine stellvertretende Implementierung liefern.
- Bei Unsicherheit mit einer kleinen Hinweistrecke helfen und den Lernenden danach wieder selbst arbeiten lassen.
- Vor längeren Erklärungen kurz erfragen, was der Lernende bereits denkt oder versucht hat.
- Der Lernende wollte keinen eigenen Papierkampf durchführen. Die KI hat einen vollständigen Kampf anhand der Regeln nachvollzogen und M0 daraufhin als abgeschlossen bewertet.
- Antworten sorgfältig auf Rechtschreibung und Grammatik prüfen.

## Repository-Stand

- Git-Remote: `https://github.com/AlRowne/roboforge-arena.git`
- Go-Modul: `github.com/AlRowne/roboforge-arena`
- Go-Version in `go.mod`: `1.26.3`
- Aktueller Branch: `main`
- Beim Erstellen dieser Übergabe waren die bisherigen Änderungen committed und der Arbeitsbaum sauber.
- Vorhandene fachliche Hauptquelle: `docs/combat-rules.md`

## M0 – abgeschlossen

Die Kampfregeln wurden gemeinsam entworfen, bereinigt und als Regelversion `v1` eingefroren. `docs/combat-rules.md` ist die verbindliche Quelle; bei Unklarheiten nicht stillschweigend davon abweichen.

Wesentliche Entscheidungen:

- Ein Bot besitzt genau ein Chassis und null bis zwei Module; Module dürfen doppelt gewählt werden.
- HP und AP starten am Maximum, EP bei abgerundeten 50 Prozent.
- Ab Runde 2 regenerieren beide Bots am Rundenanfang 10 EP.
- Eine Runde hat eine eindeutige Initiative und höchstens eine Aktion pro lebendem Bot.
- Aktionen: Angriff, Verteidigung und starker Angriff.
- Verteidigung würfelt einen Buff von 10 bis 20. Er gilt bis zum nächsten Angriff, wird danach verbraucht und kann nicht gestapelt werden.
- Schaden wird nach Verteidigung zuerst auf AP und als Überschuss auf HP angewendet.
- Strategien bestehen aus einer bis drei lückenlos priorisierten Bedingung-Aktion-Regeln; die erste erfüllte und ausführbare Regel gewinnt, sonst erfolgt ein normaler Angriff.
- Taktisch schlechte Strategien sind bewusst erlaubt.
- Nach Runde 10 endet ein noch laufender Kampf unentschieden.
- Snapshot von Bot A und Bot B, Bot-Reihenfolge, Regelversion und Seed bestimmen Ergebnis und Event-Log vollständig.
- Jedes Match besitzt eine eigene, mit seinem Seed initialisierte Zufallsquelle. Andere Matches, Worker und die Systemzeit dürfen den Ablauf nicht beeinflussen.

Die KI führte einen vollständigen Probekampf bis Runde 10 durch. Dabei wurden EP-Regeneration, starker Angriff, Verteidigung, AP-Überlauf, Prioritätswechsel, eine erfüllte aber nicht ausführbare Regel, Standardangriff und K. o. vor der zweiten Aktion erfolgreich nachvollzogen. Es war keine zusätzliche fachliche Entscheidung nötig. M0 gilt daher als bestanden.

## M1 – aktueller Lernstand

Ziel ist eine reine, deterministische Simulation ohne HTTP, PostgreSQL, Systemzeit oder globale Zufallsquelle.

Der Lernende hat die vier fachlichen Datengruppen verstanden:

1. **Simulationseingabe:** Seed, Regelversion sowie geordnete unveränderliche Snapshots beider Bots mit Ausrüstung, Maximalwerten, Speed und Strategie.
2. **Interner Kampfzustand:** aktuelle Runde, Initiative der aktuellen Runde, aktuelle HP/AP/EP und aktive Verteidigungsbuffs.
3. **Kampfereignisse:** unveränderliche Aufzeichnungen bereits abgeschlossener Vorgänge.
4. **Simulationsergebnis:** Sieger oder Unentschieden mit Grund, Endrunde, Endzustände beider Bots und geordnetes Event-Log; ungültige Eingaben führen stattdessen zu einem Fehler.

Folgende Event-Granularität wurde entschieden:

- Ein Rundenstart-Event pro Runde dokumentiert Rundennummer, tatsächliche EP-Änderung beider Bots und Initiative samt Grund.
- Ein zusammengefasstes Aktions-Event pro ausgeführter Aktion dokumentiert Runde und Aktionsposition, Akteur und Ziel, ausgewählte Aktion, auslösende Strategiepriorität oder Default, Zufallswert, angewandte Verteidigung, effektiven Schaden, Ressourcenänderungen, resultierenden Zustand und K.-o.-Status.

Der Lernende versteht den Persistenzfluss:

> Simulation sammelt Events geordnet im Arbeitsspeicher und gibt ein Ergebnis zurück → Worker erhält das Ergebnis → Worker speichert Ergebnis, Events und finalen Match-Status gemeinsam in einer Datenbanktransaktion.

Die Simulation kennt weder Worker noch Datenbank. Eine geordnete Folge ist für das Event-Log passender als eine Map.

## Unmittelbar nächster Schritt

Der Lernende möchte jetzt in den Go-Code einsteigen. `go.mod` wurde bereits vom Lernenden erstellt; diesen Schritt nicht wiederholen.

Der vereinbarte kleinste Implementierungsschritt lautet:

> Aus einer gültigen Botkonfiguration wird ohne Seiteneffekte der korrekte Anfangszustand erzeugt.

Empfohlener Arbeitsbereich ist ein reines internes Simulationspaket, voraussichtlich `internal/simulation`. Die endgültigen Go-Typen, Namen und Paketentscheidungen trifft der Lernende.

Erster fachlicher Testfall:

- Eingabe: Chassis A mit Modul A und Modul C.
- Erwartete Maximalwerte: 140 HP, 100 AP, 100 EP und 120 Speed.
- Erwarteter Startzustand: 140 aktuelle HP, 100 aktuelle AP, 50 aktuelle EP und kein Verteidigungsbuff.
- Die übergebene Konfiguration bleibt unverändert.

Der Lernende soll zuerst selbst den Test und anschließend nur die minimale Implementierung dafür schreiben. Bei der nächsten Unterhaltung sollte die KI ihn seinen Ansatz beziehungsweise den entstandenen Code erklären lassen, dann gegen diesen Testfall und die Regeln prüfen. Keine fertige Test- oder Implementierungslösung vorgeben.

## Noch offene Entscheidungen in M1

Diese Punkte wurden bewusst noch nicht festgelegt und sollen schrittweise vom Lernenden entschieden werden:

- konkrete Go-Typen und Bezeichner für Konfiguration, Zustand, Events und Ergebnis;
- genaue Paketstruktur der Simulation;
- konkreter reproduzierbarer Zufallszahlengenerator und dessen Kapselung;
- Repräsentation der beiden Event-Arten in Go;
- Validierungsreihenfolge und Fehlerdarstellung;
- Aufbau der Kampfschleife und Aufteilung in kleine, testbare Schritte.

Beim Einstieg nicht die vollständige Simulation auf einmal entwerfen. Zuerst den Anfangszustand als kleinen grünen Test abschließen, dann den nächsten vertikalen Schritt gemeinsam bestimmen.
