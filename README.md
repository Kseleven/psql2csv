# Psql2csv
a postgresql tool,support export psql data to csv and import csv back to psql.

### Feature
support dataType:
* text
* bigint
* boolean
* text[]
* numeric
* inet
* timestamp

### Usage
install:
```shell
git clone https://github.com/Kseleven/psql2csv.git
go run cmd/pg2csv.go --help
```

configuration file is below examples path,copy and modify it,then execute command below:
```shell
go run cmd/pg2csv.go -f {yourpath}/examples/example.json
```