# GIVE API 

1. Edit `src/config.yaml` file

2. Compile REST API from `src` folder

* Install [Go](https://golang.org/doc/install#install)

* Run `make get_libs` to install required libraries

* Run `make test` to run unit tests

* Run `make` to compile api

3. Launch API

`./give_api`

Output:
```
./give_api 
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /api/v1/give/kids         --> main.CreateKids (5 handlers)
[GIN-debug] PUT    /api/v1/give/kids         --> main.UpdateKids (5 handlers)
[GIN-debug] GET    /ping                     --> main.restEngine.func1 (4 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

4. Available API Endpoints

#### `Authentication`

   All endpoints require an auth header like this : 
   `GIVEAPIToken: c7ee388b68b8766f51631da056b2bbd9`
    (value for token can be changed on `config.yaml`)

#### `POST   /api/v1/give/kids` 

   *Parameters*: JSON object containing kid data, e.g: `{"id":"0x9AbAB02EcBe8A917C266681B37d1f45f56191bDb", "name":"Danny Wood", "date_of_birth":"2010-11-10", "parents_emails":["pa@gmail.com","ma@gmail.com"], "school_name":"1st School of Hawaii", "id_tag_name":"Big Danny"}`

   *Returns*: JSON kid object.  e.g. :  `{"kids":[{"id":"0x9AbAB02EcBe8A917C266681B37d1f45f56191bDb","name":"Danny Wood","date_of_birth":"2010-11-10","parents_emails":["pa@gmail.com","ma@gmail.com"],"students_photo":null,"school_name":"1st School of Hawaii","id_tag_name":"Big Danny"}]}`

   Curl example:
```
curl -s -H "Content-Type: application/json" -H "GIVEAPIToken: $TOKEN" -X POST -d '{"id":"0x9AbAB02EcBe8A917C266681B37d1f45f56191bDb", "name":"Danny Wood", "date_of_birth":"2010-11-10", "parents_emails":["pa@gmail.com","ma@gmail.com"], "school_name":"1st School of Hawaii", "id_tag_name":"Big Danny"}' 'http://localhost:8080/api/v1/give/kids'
```
