# POC gRPC - Microservices - Containers
# Services
## User data service
Manages the creation, authentication and basic information of users
### Envoreontment variables
- udsConnStr: onnection string used to connect with a mysql database, must be in format: `[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]` 

    More examples of valid format found in: https://github.com/go-sql-driver/mysql
- udsListenAddr: address in which the service will be listening
### Exposed methods
- CreateUser
- GetUser

### CreateUser request
```
{
    "authToken": "incididunt sunt nostrud elit",
    "user": {
        "age": 59,
        "name": "Teresa Fajardo",
        "parentId": 0,
        "pwdHash": "some other password"
    }
}
```
### CreateUser response
```
{
    "id": "3"
}
```
### GetUser request
```
{
    "authToken": "pariatur",
    "id": "3"
}
```
### GetUser response
```
{
    "user": {
        "id": "3",
        "name": "teresa fajardo",
        "pwdHash": "some other password",
        "age": 59
    },
    "status": {}
}
```
__*\*Currently authentication is in development, so auth tokens can be any non empty string*__
