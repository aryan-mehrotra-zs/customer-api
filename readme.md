# customer-api
This project allows to perform crud operations.

## Steps to setup server

### To setup database container
```azure
docker run  --name customer-api -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=organisation -p 3306:3306 -d mysql:latest
```
```azure
docker exec -it sample-api mysql -u root -ppassword organisation
```

### MySQL Commands

```azure
use organisation

CREATE TABLE customers (
                            ID int NOT NULL AUTO_INCREMENT,,
                            name varchar(255) NOT NULL,
                            phone_no int NOT NULL,
                            address varchar(255) NOT NULL,
                            PRIMARY KEY (ID)
);
```

### To start the server

`go run main.go`
