# RoboForge Arena – Kampfregeln v1

## 1. Roboterkonfiguration

Ein Roboter (Bot) besteht aus genau einem Chassis und höchstens zwei Modulen. Das Chassis legt die Basiswerte fest. Module erhöhen einzelne Basiswerte; dasselbe Modul darf mehrfach gewählt werden.

Ein Bot besitzt folgende Werte:

- Lebenspunkte (HP)
- Schildpunkte (AP)
- Energiepunkte (EP)
- Geschwindigkeit (Speed)

### Chassis

Die Werte eines Chassis sind die vollständigen Basiswerte des Bots.

| Wert | Chassis A | Chassis B |
|---|---:|---:|
| HP | 120 | 100 |
| AP | 100 | 120 |
| EP | 100 | 120 |
| Speed | 100 | 100 |

### Module

| Wert | Modul A | Modul B | Modul C |
|---|---:|---:|---:|
| Speed | +20 | 0 | 0 |
| AP | 0 | +20 | 0 |
| HP | 0 | 0 | +20 |

Die Maximalwerte eines Bots ergeben sich aus seinen Chassiswerten und allen Modulboni.

## 2. Ausgangszustand

Zu Beginn eines Kampfes gilt:

- HP und AP entsprechen ihren jeweiligen Maximalwerten.
- EP entsprechen 50 Prozent des Maximalwertes und werden bei einem nicht ganzzahligen Ergebnis abgerundet.
- Es ist kein Verteidigungsbuff aktiv.

HP, AP und EP können niemals unter 0 fallen. Sie können ihre jeweiligen Maximalwerte nicht überschreiten. HP und AP regenerieren sich während eines Kampfes nicht.

Ein Bot ist außer Gefecht, sobald seine HP 0 erreichen.

## 3. Kampfablauf

Ein Kampf dauert höchstens zehn Runden. Jede Runde läuft in dieser Reihenfolge ab:

1. Ab Runde 2 regenerieren beide Bots gleichzeitig 10 EP. Dabei kann der jeweilige EP-Maximalwert nicht überschritten werden.
2. Die Initiative wird bestimmt. Der Bot mit dem höheren Speedwert erhält die Initiative. Bei gleichem Speed haben beide Bots dieselbe Chance, die Initiative für diese Runde zu erhalten.
3. Der Bot mit der Initiative wertet seine Strategie anhand des aktuellen Kampfzustands aus und führt die ausgewählte Aktion sofort vollständig aus.
4. Fällt dabei ein Bot auf 0 HP, endet der Kampf sofort.
5. Andernfalls wertet der zweite Bot seine Strategie anhand des nun veränderten Kampfzustands aus und führt seine Aktion sofort vollständig aus.
6. Fällt dabei ein Bot auf 0 HP, endet der Kampf sofort. Andernfalls endet die Runde.
7. Sind nach Abschluss der zehnten Runde beide Bots noch kampffähig, endet der Kampf unentschieden.

Jeder kampffähige Bot kann höchstens eine Aktion pro Runde ausführen. Ein außer Gefecht gesetzter Bot führt keine weitere Aktion aus.

## 4. Aktionen

Alle angegebenen Zahlenbereiche schließen den unteren und oberen Grenzwert ein. Jeder ganzzahlige Wert innerhalb eines Bereichs besitzt dieselbe Wahrscheinlichkeit.

| Aktion | EP-Kosten | Effekt | Ausführbar, wenn |
|---|---:|---|---|
| Angriff | 0 | Verursacht 20–30 Schaden | immer |
| Verteidigung | 0 | Erzeugt einen Verteidigungsbuff mit einem Wert von 10–20 | kein Verteidigungsbuff aktiv ist |
| Starker Angriff | 20 | Verursacht 35–45 Schaden | mindestens 20 EP vorhanden sind |

Die Kosten eines starken Angriffs werden bei seiner Ausführung sofort von den aktuellen EP abgezogen. Erst danach wird sein Schaden bestimmt.

### Verteidigung

Bei der Ausführung von Verteidigung wird der Wert des Buffs einmal zufällig bestimmt. Der Buff bleibt bis zum nächsten eingehenden normalen oder starken Angriff aktiv.

Bei diesem Angriff wird der Verteidigungswert vom ausgewürfelten Schaden abgezogen. Der verbleibende Schaden kann nicht negativ werden. Anschließend wird der Verteidigungsbuff auch dann entfernt, wenn er den gesamten Schaden verhindert hat.

Solange ein Verteidigungsbuff aktiv ist, kann der betreffende Bot nicht erneut Verteidigung ausführen.

### Schadensreihenfolge

Ein normaler oder starker Angriff wird folgendermaßen verarbeitet:

1. Der Schaden wird innerhalb des Effektbereichs der Aktion ausgewürfelt.
2. Bei einem aktiven Verteidigungsbuff wird dessen Wert vom Schaden abgezogen und der Buff anschließend entfernt.
3. Der verbleibende Schaden wird zuerst von den AP abgezogen.
4. Übersteigt der Schaden die aktuellen AP, fallen die AP auf 0 und der exakte Überschuss wird von den HP abgezogen.
5. Fallen die HP durch den Schaden unter 0, werden sie auf 0 begrenzt.

Beispiel: Ein Bot mit 5 AP wird nach Berücksichtigung eines möglichen Verteidigungsbuffs von 20 Schaden getroffen. Seine AP fallen auf 0 und seine HP werden um 15 reduziert.

## 5. Strategieregeln

Eine Strategie besteht aus mindestens einer und höchstens drei priorisierten Regeln. Jede Regel kombiniert frei eine der erlaubten Bedingungen mit einer der erlaubten Aktionen.

### Erlaubte Bedingungen

1. Die aktuellen HP des Gegners betragen höchstens 25 Prozent seiner maximalen HP einschließlich Modulboni.
2. Die eigenen AP sind 0.
3. Die eigenen EP reichen für einen starken Angriff, also sind mindestens 20 EP vorhanden.

### Erlaubte Aktionen

- Angriff
- Verteidigung
- Starker Angriff

### Prioritäten und Auswahl

Die Prioritäten beginnen bei 1, sind eindeutig und lückenlos. Priorität 1 ist die höchste Priorität. Bei zwei Regeln sind daher die Prioritäten 1 und 2 zulässig, bei drei Regeln die Prioritäten 1, 2 und 3.

Unmittelbar vor seinem Zug wählt ein Bot seine Aktion wie folgt:

1. Die Regeln werden von der höchsten zur niedrigsten Priorität geprüft.
2. Eine Regel ist anwendbar, wenn ihre Bedingung erfüllt und ihre Aktion ausführbar ist.
3. Die erste anwendbare Regel bestimmt die Aktion; weitere Regeln werden nicht mehr geprüft.
4. Ist eine Bedingung nicht erfüllt oder ihre Aktion nicht ausführbar, wird die nächste Regel geprüft.
5. Ist keine Regel anwendbar, führt der Bot einen normalen Angriff aus.

Eine taktisch schlechte Strategie bleibt gültig. Das Spiel verhindert beispielsweise nicht, dass ein Bot wegen seiner Prioritäten wiederholt verteidigt und dadurch möglicherweise verliert oder nach zehn Runden ein Unentschieden erreicht.

## 6. Zufall und Determinismus

Jedes Match erhält beim Erstellen genau einen zufällig erzeugten Seed. Dieser Seed bleibt dem Match dauerhaft zugeordnet.

Jede Matchsimulation besitzt eine eigene Zufallsquelle, die mit dem Seed dieses Matches initialisiert wird. Zufallsquellen anderer Matches, die Anzahl der Worker, die Uhrzeit und die Verarbeitungsgeschwindigkeit dürfen den Kampf nicht beeinflussen.

Eine Zufallszahl wird ausschließlich für folgende Vorgänge benötigt:

1. Initiative einer Runde bei gleichem Speed
2. Schaden eines normalen Angriffs
3. Schaden eines starken Angriffs
4. Wert eines Verteidigungsbuffs

Die Zufallswerte werden in der zeitlichen Reihenfolge der Kampfereignisse gezogen. Eine nicht ausgeführte Aktion verbraucht keinen Zufallswert.

Die vollständigen Eingaben einer Simulation sind:

- der unveränderliche Snapshot von Bot A,
- der unveränderliche Snapshot von Bot B,
- die eindeutige Reihenfolge beziehungsweise Identität beider Bots,
- die Regelversion `v1`,
- der Match-Seed.

Dieselben vollständigen Eingaben müssen dasselbe Kampfergebnis und dasselbe Event-Log in derselben Reihenfolge erzeugen. Dazu gehören insbesondere Initiative, ausgewählte Aktionen, gewürfelte Werte, Ressourcenstände nach jeder Aktion und der Sieger.

Wird ein Kampf durch einen Workerfehler unterbrochen und wurde kein vollständiges Ergebnis gespeichert, wird er ab Runde 1 mit denselben Eingaben neu simuliert. Dadurch entsteht erneut derselbe vollständige Kampf.

Änderungen am Kampfverhalten oder Balancing benötigen eine neue Regelversion. Dazu gehören beispielsweise veränderte Effektbereiche, Energiekosten, Bedingungen, Regeneration oder Rundengrenzen. Reine Rechtschreibkorrekturen oder betriebliche Änderungen wie eine andere Anzahl von Workern benötigen keine neue Regelversion.

## 7. Kampfende

Ein Bot verliert sofort, wenn seine HP 0 erreichen. Der andere Bot gewinnt und der Kampf endet vor jeder weiteren Aktion.

Sind nach der vollständig ausgeführten zehnten Runde beide Bots noch kampffähig, endet der Kampf unentschieden.

Mit den Aktionen der Regelversion `v1` können nicht beide Bots durch dieselbe Aktion gleichzeitig außer Gefecht gesetzt werden.

## 8. Ungültige Konfigurationen

Eine Botkonfiguration ist ungültig, wenn mindestens einer der folgenden Fälle eintritt:

- Es wurde kein Chassis oder mehr als ein Chassis gewählt.
- Das Chassis ist weder Chassis A noch Chassis B.
- Es wurden mehr als zwei Module gewählt.
- Ein gewähltes Modul ist weder Modul A, Modul B noch Modul C.
- Die Strategie enthält weniger als eine oder mehr als drei Regeln.
- Eine Regel enthält eine unbekannte Bedingung oder Aktion.
- Eine Regel besitzt keine Priorität.
- Eine Priorität kommt mehrfach vor.
- Die Prioritäten beginnen nicht bei 1 oder enthalten eine Lücke.
- Dieselbe Kombination aus Bedingung und Aktion kommt mehrfach vor.

Null, ein oder zwei Module sind gültig. Module dürfen mehrfach gewählt werden. Eine taktisch schlechte, aber strukturell gültige Strategie wird nicht abgelehnt.

## 9. Zentrale Invarianten

- Gleiche vollständige Eingaben erzeugen dasselbe Ergebnis und dasselbe Event-Log.
- HP, AP und EP liegen immer zwischen 0 und ihrem jeweiligen Maximum.
- Ein außer Gefecht gesetzter Bot führt keine weitere Aktion aus.
- Jeder kampffähige Bot führt höchstens eine Aktion pro Runde aus.
- Jede ausgeführte Aktion war unmittelbar vor ihrer Ausführung erlaubt.
- Ein Kampf endet spätestens nach der zehnten Runde.
