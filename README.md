# bede
Just a dictionary

# Notes
This code assumes that the folder to parse contains a relatively small number of files and those files are not too big. Therefore, we generate the whole dictionary in memory before writing it to disk.

For a large number of file, it would be preferrable to stream the dictionaries to disk as you walk the source folder instead of loading in it memory.

# Tests
You can run `go test` to run the tests

# CLI
The CLI is located in `cmd`. To use it:
```
cd cmd
go build -o bede 
./bede --help
./bede -source "../testdata" -client "client.xml" -server "server.xml"
```
