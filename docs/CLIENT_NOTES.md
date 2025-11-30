# Klijentske napomene / zahtjevi (aktualno)

- **Podaci pacijenta**: Ime, Prezime, Broj (telefon), Adresa, Datum rođenja, Spol.
- **Struktura zapisa posjete (anamneza)**:
  - Anamneza (bilješka sesije)
  - Dijagnoza
  - Terapija
  - Ostale informacije (npr. drugi dolazak, stanje)
  - Razlog posjete (prikazati na PDF-u umjesto generičke evaluacije)
  - Datumi dolazaka (npr. 20.10 prva, 30.10 druga). Zadnji PDF treba sadržavati sve prethodne.
- **PDF**:
  - PDF se stalno ažurira (upsert): zadnji PDF sadrži sve prethodne anamneze/dijagnoze/terapije, uz zadnju posjetu/razlog.
  - Ako je više ciklusa (pauza → novi ciklus), generirati novi dokument, ali prethodni ostaju.
  - Kod nove tegobe (npr. došao zbog koljena, kasnije leđa), PDF se ažurira adekvatno, ali anamneza/dijagnoza/terapija starih posjeta se ne mijenjaju.
- **Pretraga/kategorizacija**:
  - Potencijalna pretraga po nalazima/tegobama (npr. leđa, koljeno, kuk).
  - Kategorizacija po čemu je riječ (tip tegobe/razlog posjete).
- **Ponašanje**:
  - Sve prethodno što je pisano ostaje; prikazuje se zadnji put i razlog posjete.
  - Ako ima 5 posjeta za 5 različitih stvari, ispisati razlog posjete po posjeti.
