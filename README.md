# Psql2csv
A postgresql tool,support export psql data to csv and import csv back to psql.

### Feature
Support dataType:
* text
* bigint
* boolean
* text[]
* numeric
* inet
* timestamp

### Usage
clone source:
```shell
git clone https://github.com/Kseleven/psql2csv.git
go run cmd/pg2csv.go --help
```

install
```shell
go install github.com/Kseleven/psql2csv@latest
pg2csv --help
```

export and import tables:
1. configuration file is below examples path,copy and modify it,then execute command below:
    ```shell
    go run cmd/pg2csv.go -f {yourpath}/examples/example.json
    ```

2. You wil get table.csv files in "exportPath" specified in example.json.
3. Change parameter "action" in example.json to "import", and modify new database config
4. Execute command to import csv to new database
    ```shell
        go run cmd/pg2csv.go -f {yourpath}/examples/example.json
    ```
5. If new database column is diff with csv header, it will print different column. 
Then you should add or delete csv header to keep same with new database column.
6. After modify csv header and data, execute import command again.  