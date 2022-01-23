#Customer-api

## Setup Docker
###Linux
#### Install Docker        `` sudo apt install docker``
#### Download MySql Image  ``docker run  --name customer-api -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=organisation -p 3306:3306 -d mysql:latest``
#### Here ,                ``Username: root Password: password Database: organisation``
#### Execute Docker Image  ``docker exec -it customer-api mysql -u root -ppassword organisation``
#### Switch to Database    ``USE organisation``
#### Create Table          
```
CREATE TABLE table_name (
    ID int NOT NULL AUTO_INCREMENT,,
    name varchar(255) NOT NULL,
    address varchar(255) NOT NULL,
    phone_no int NOT NULL,
    PRIMARY KEY (ID)
); 
```




// SERVICE DELETE - ERROR
// ELSE DATA,ERROR

// STORE CREATE - ID(WHEN AUTO GENERAATED IN DB),ERROR
// DELETE - ERROR
// UPDATE - ERROR
// GET - DATA,ERROR

















#### Remove Sudo (Optional) ``sudo setfacl -m user:$USER:rw /var/run/docker.sock``


