# Roboterzustand
- Ein Roboter besitzt Lebenspunkte (HP), Schildpunkte (AP) und Energie (EP) und einen Speedwert.
- Wenn die HP auf 0 sind ist er ausser Gefecht.
- Zuerst muessen die AP auf 0 sein, bevor HP abgezogen werden.
- Energie wird fuer verschiedene Aktionen verwendet. Derzeit nur starker Angriff. Sie startet bei 50% des Maximalwertes.
- Es gibt einen normalen Angriff und Verteidigen, was keine Energie kostet.
- Energie regeneriert sich am Anfang der Runde um 10. Der Maximalwert kann nicht ueberschritten werden.
- HP regenieren sich nicht.
- Alle werte werden durch die Roboterkonfiguration beeinflusst.
1. Das Chassis gibt die Basiswerte vor.
2. Das Modul erhoeht diese Werte.
- Ein Roboter hat ein Chassis und hoechstens zwei Module.

# Kampfablauf
- Jeder Robo startet mit maximalen HP und AP und 50% EP, abgerundet.
- Der Roboter mit dem hoehren Speedwert startet. Bei Gleichstand entscheidet der Zufall.
- Am Start jeder Runde ab Runde 2 regenerieren sich EP einmalig um 10 fuer beide Robos gleichzeitig. Die Werte gehen niemals ueber den Maximalwert.
- Jeder Roboter entscheided nacheinander seine Aktionen.
- Die Aktion des Robos mit der hoeheren Initiative wird zuerst durchgefuehrt. Das ist entweder der mit dem hoeheren Speed-Stat, oder bei Speedgleichheit der, der den Zufallsroll fuer diese Runde gewonnen hat. Diese Aktion wird sofort durchgefuehrt, erst abgeschlossener Berechnung dieser Aktion wird die des anderen Robos durchgefuehrt (falls er noch mehr als 0 HP hat).
- Schaden, der ueber den AP-Wert geht, wird den HP abgezogen.
- Die Runde endet, wenn beide Aktionen durchgefuehrt wurden, oder wenn ein Robo auf 0 HP faellt. Dann Endet auch der Kampf.

# Aktionen
- Es gibt Aktionen wie Angriff, Verteidigung. Diese brauchen keine Energie. Dann gibt es spezielle Aktionen, die Energie brauchen. Das ist derzeit Starker Angriff. 
- Basic Aktionen kosten nichts. Das sind Angriff und Verteidigung. Spezielle Aktionen kosten Energie, starker Angriff 20 EP.
- Die EP werden sofort bei Ausfuehrung der Aktion abgezofen.
- Die Verteidigungsaktion schwaecht den naechsten Angriff der Gegners ab. Dabei wird der Verteidigungswert dem Angriffswert abgezogen. Dann wird der Verteidigungsbuff entfernt. Der Schade kann nicht ins negative gehen. Es kann nur verteidigt werden, wenn kein Verteidigungsbuff aufrecht ist.
- Es kann nur eine Aktion pro Runde ausgewaehlt werden.

# Strategieregeln
- Hier kann der Spiele vorgeben, welche Aktionen priorisiert werden, also eher Angriff oder eher Verteidigung. Diese werden dann bevorzugt ausgefuehrt. 
- Es werden 3 Bedingungen vorgegeben:
1. Gegnerische HP <= 25 %
2. Eigene AP sind leer
3. Energie reicht fuer starken Angriff.
- Diesen Bedingungen sollen Aktionen zugewiesen werden:
1. Starker Angriff
2. Verteidigung
3. Angriff
- Es sollen mindesten 1, maximal 3 Regeln erstellt werden koennen.
- Diese Regeln koennen frei kombiniert werden. Eine Bedinung muss mit einer Aktion kombiniert werden. Es darf dieselbe Kombination nicht doppelt vorkommen.
- Diese Regeln sollen soll eine Prioritaet zugewiesen werden.
- Der Auswahlauflauf funktioniert so:
1. Hoechste Prioritat.
2. Bedingung erfuellt? -> Machbarkeitscheck -> ja: ausfuehren, nein: naechste Prioritaet -> Bedinungs erfuellt? -> etc.
- Machbarkeitscheck bedeuted:
1. Angriff ist immer ausfuehbar.
2. Starker Angriff nur mit mindestens 20 EP.
3. Verteidigung ist nur ohne aktiven Verteidigungsbuff ausfuehrbar.
- Wenn keine Regel anwendbar ist, soll ein Angriff ausgefuehrt werden.

# Seed
- Die Simulation besitzt ihre eigene Zufallsquelle, die unabhaenging von den Seeds der anderen Matches ist.
- Fuer jedes Match wird beim Erstellen ein zufaelliger Seed erzeugt, der genau diesem Match zugeordnet bleibt, auch wenn es zB. durch einen Fehler im Worker unterbrochen wurde. Dadurch kann dieses Match neu simuliert werden und wird wieder in derselben Weise ablaufen.
- Eine Zufallszahl wird benoetigt bei:
1. Initiative bei Speed-Gleichheit
2. Angriffsschaden
3. starkem Angriffsschaden
4. Verteidigungswert
- Die Zufallswerte werden in zeitlicher Reihenfolge der Kampfereignisse gezogen

# Zufall
- Bei gleichem Speedstat entscheidet der Zufall, welcher Robo zuerst handeln kann. Dieser Zufallscheck wird dann vor jeder neuen Runde durchgefuehrt.
- Der Schaden unterliegt einem Bereich. Er ist inklusiv gemeint und steht bei den Aktionen.
- Fuer die Verteidigung gilt dasselbe.

# Kampfende
- Der Roboter, der zuerst 0 HP hat, verliert. Die Aktion dieses Robos wird dann nicht mehr ausgefuehrt.
- Nach dem Ende der 10. Runde ist unentschieden, sofern keiner der Robos 0 HP hat.

# Ungueltige Konfigurationen
- Es kann nur ein Chassis gewaehlt werden. A oder B.
- Es koennen maximal 2 Module gewaehlt werden. Module koennen mehrfach gewaehlt werden. Es stehen Module A, B und C zur Auswahl.
- Es muss eine Strategie gewaehlt werden.
- Es duerfen nur vorgegebene Aktionen oder Strategien verwendet werden.

# Equipment
| Chassis | A   | B   |
|---------|-----|-----|
| HP      | 120 | 100  |
| AP      | 100  | 120 |
| Energie | 100 | 120 |
| Speed   | 100 | 100 |

| Modul   | A   | B   | C   |
|---------|-----|-----|-----|
| Speed   | +20 | +0  | +0  |
| AP | +0  | +20 | +0  |
| HP | +0  | +0  | +20 |

# Aktionen
| Aktion             | Energiekosten | Effektbereich | Bedingungen       |
|--------------------|---------------|-----------------|-------------------|
| Angriff            | 0             | 20-30           |                   |
| Verteidigung       | 0             | 10-20          |                   |
| Starker Angriff    | 20            | 35-45           |                   |
