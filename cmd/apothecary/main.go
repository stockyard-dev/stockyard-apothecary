package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-apothecary/internal/server";"github.com/stockyard-dev/stockyard-apothecary/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="10200"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./apothecary-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("apothecary: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Apothecary — health and medication tracker\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("apothecary: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
