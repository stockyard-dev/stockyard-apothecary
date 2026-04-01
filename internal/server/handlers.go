package server
import("encoding/json";"net/http";"github.com/stockyard-dev/stockyard-apothecary/internal/store")
func(s *Server)handleListMeds(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListMedications();if list==nil{list=[]store.Medication{}};writeJSON(w,200,list)}
func(s *Server)handleCreateMed(w http.ResponseWriter,r *http.Request){var m store.Medication;json.NewDecoder(r.Body).Decode(&m);if m.Name==""{writeError(w,400,"name required");return};m.Active=true;s.db.CreateMedication(&m);writeJSON(w,201,m)}
func(s *Server)handleLogEntry(w http.ResponseWriter,r *http.Request){var e store.Entry;json.NewDecoder(r.Body).Decode(&e);if e.Type==""||e.Description==""{writeError(w,400,"type and description required");return};s.db.LogEntry(&e);writeJSON(w,201,e)}
func(s *Server)handleListEntries(w http.ResponseWriter,r *http.Request){typ:=r.URL.Query().Get("type");list,_:=s.db.ListEntries(typ);if list==nil{list=[]store.Entry{}};writeJSON(w,200,list)}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
