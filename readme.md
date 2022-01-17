#Customer-api

## Setup Docker
###Linux
#### Install Docker        `` sudo apt install docker``
#### Download MySql Image  ``docker run  --name customer-api -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=organisation -p 3306:3306 -d mysql:latest``
#### Here ,                ``Username: root Password: password Database: organisation``
#### Execute Docker Image  ``docker exec -it sample-api mysql -u root -ppassword organisation``
#### Switch to Database    ``USE organisation``
#### Create Table          
```
CREATE TABLE table_name (
    ID int NOT NULL,
    Name varchar(255) NOT NULL,
    PhoneNo int NOT NULL,
    Address varchar(255) NOT NULL,
    PRIMARY KEY (ID)
); 
```
























#### Remove Sudo (Optional) ``sudo setfacl -m user:$USER:rw /var/run/docker.sock``


