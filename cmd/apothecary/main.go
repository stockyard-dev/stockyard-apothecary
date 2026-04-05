package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-apothecary/internal/server";"github.com/stockyard-dev/stockyard-apothecary/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./apothecary-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("apothecary: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Apothecary — Self-hosted health and medication tracker\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n",port,port)
log.Printf("apothecary: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
