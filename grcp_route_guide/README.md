Route Guide gRPC
=====

```bash
#To run the server
go run server/server.go -json_db_file testdata/route_guide_db.json

#To run the client
go run client/client.go

```

-----
This example show the four service methods supported by gRPC

* Unary
* Server-side streaming
* Client-side streaming
* Bidirectional