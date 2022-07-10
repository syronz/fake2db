# fake2db
Create fake data and bulk insert it into the database


## sample table 
Below table in mysql is used as a sample

```sql
CREATE TABLE students (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    gender ENUM('male', 'female', 'other'),
    code INT,
    dob DATE,
    address TEXT,
    created_at DATETIME,
    PRIMARY KEY (id)
);
```
