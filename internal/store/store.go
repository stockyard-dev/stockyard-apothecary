package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ db *sql.DB }
type Medication struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Dosage     string `json:"dosage"`
	Frequency  string `json:"frequency"`
	Prescriber string `json:"prescriber"`
	Pharmacy   string `json:"pharmacy"`
	RefillDate string `json:"refill_date"`
	Notes      string `json:"notes"`
	Active     int    `json:"active"`
	CreatedAt  string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "apothecary.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS medications(id TEXT PRIMARY KEY,name TEXT NOT NULL,dosage TEXT DEFAULT '',frequency TEXT DEFAULT '',prescriber TEXT DEFAULT '',pharmacy TEXT DEFAULT '',refill_date TEXT DEFAULT '',notes TEXT DEFAULT '',active INTEGER DEFAULT 1,created_at TEXT DEFAULT(datetime('now')))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
	resource TEXT NOT NULL,
	record_id TEXT NOT NULL,
	data TEXT NOT NULL DEFAULT '{}',
	PRIMARY KEY(resource, record_id)
)`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string        { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string          { return time.Now().UTC().Format(time.RFC3339) }
func (d *DB) Create(e *Medication) error {
	e.ID = genID()
	e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO medications(id,name,dosage,frequency,prescriber,pharmacy,refill_date,notes,active,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`, e.ID, e.Name, e.Dosage, e.Frequency, e.Prescriber, e.Pharmacy, e.RefillDate, e.Notes, e.Active, e.CreatedAt)
	return err
}
func (d *DB) Get(id string) *Medication {
	var e Medication
	if d.db.QueryRow(`SELECT id,name,dosage,frequency,prescriber,pharmacy,refill_date,notes,active,created_at FROM medications WHERE id=?`, id).Scan(&e.ID, &e.Name, &e.Dosage, &e.Frequency, &e.Prescriber, &e.Pharmacy, &e.RefillDate, &e.Notes, &e.Active, &e.CreatedAt) != nil {
		return nil
	}
	return &e
}
func (d *DB) List() []Medication {
	rows, _ := d.db.Query(`SELECT id,name,dosage,frequency,prescriber,pharmacy,refill_date,notes,active,created_at FROM medications ORDER BY created_at DESC`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Medication
	for rows.Next() {
		var e Medication
		rows.Scan(&e.ID, &e.Name, &e.Dosage, &e.Frequency, &e.Prescriber, &e.Pharmacy, &e.RefillDate, &e.Notes, &e.Active, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}
func (d *DB) Update(e *Medication) error {
	_, err := d.db.Exec(`UPDATE medications SET name=?,dosage=?,frequency=?,prescriber=?,pharmacy=?,refill_date=?,notes=?,active=? WHERE id=?`, e.Name, e.Dosage, e.Frequency, e.Prescriber, e.Pharmacy, e.RefillDate, e.Notes, e.Active, e.ID)
	return err
}
func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM medications WHERE id=?`, id)
	return err
}
func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM medications`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []Medication {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ?)"
		args = append(args, "%"+q+"%")
	}
	if v, ok := filters["active"]; ok && v != "" {
		where += " AND active=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(`SELECT id,name,dosage,frequency,prescriber,pharmacy,refill_date,notes,active,created_at FROM medications WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Medication
	for rows.Next() {
		var e Medication
		rows.Scan(&e.ID, &e.Name, &e.Dosage, &e.Frequency, &e.Prescriber, &e.Pharmacy, &e.RefillDate, &e.Notes, &e.Active, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Stats() map[string]any {
	m := map[string]any{"total": d.Count()}
	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
