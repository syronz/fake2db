# Fake2DB
[![Go Report Card](https://goreportcard.com/badge/github.com/syronz/fake2db)](https://goreportcard.com/report/github.com/syronz/fake2db)
[![Passing](https://github.com/syronz/fake2db/actions/workflows/go.yml/badge.svg)](https://github.com/syronz/fake2db/actions/workflows/go.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/syronz/fake2db)](https://pkg.go.dev/github.com/syronz/fake2db)

Create fake data and bulk insert it into the database


## Dependency
For generating fake data using below package
[https://github.com/brianvoe/gofakeit](https://github.com/brianvoe/gofakeit)

## Notes for running the app
In case the number of workers is more than chunk like below some errors maybe accused depends on mysql configuration
```toml
worker =100 
chunk_size = 1
rows = 100000
```

supposed errors
```log
[mysql] packets.go:37: unexpected EOF
bad connection
error in executing the query Error 1213: Deadlock found when trying to get lock; try restarting transaction
```

but if chunk size is big enough it is possible to have 100 workers in the same time for instance like below
```toml
worker =100
chunk_size = 3000
```

```sql
CREATE TABLE students (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    gender ENUM('male', 'female', 'other'),
    code INT,
    dob DATE,
    address TEXT,
    created_at DATETIME,
    PRIMARY KEY (id)
);
```
