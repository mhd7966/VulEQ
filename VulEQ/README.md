# VulEQ

For running migration, run this command:
```
go get -v github.com/rubenv/sql-migrate/...
```
then:
```
sql-migrate up
```
 Swagger init changes, run in main project directory : 
```
  swag init -g cmd/api/main.go
```

  Log level values for set env as LOGLEVEL : 
```
    - panic
    - fatal
    - error (*)
    - warn/warning
    - info(*)
    - debug
    - trace
```