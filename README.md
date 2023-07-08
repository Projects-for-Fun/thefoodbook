# The Food Book

[![Github Actions](https://github.com/Projects-for-Fun/thefoodbook/actions/workflows/go.yml/badge.svg)](https://github.com/Projects-for-Fun/thefoodbook/actions/workflows/go.yml)


[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Projects-for-Fun_thefoodbook&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Projects-for-Fun_thefoodbook)


## Running locally

### Neo4j Desktop

Using Neo4j Desktop to create a db locally for the project ([link](https://neo4j.com/docs/desktop-manual/current/operations/create-dbms/)).

#### Migration 

Using golang-migrate. More info [here](https://github.com/mariamihai/neo4j-related/blob/main/readme.md#golang-migrate).

### Docker dependencies

Run `make infra-up` and run the project locally with `webservice`.

Need to set environment variables: 
```
ENVIRONMENT=

SERVICE_NAME=
SERVICE_PORT=

LOG_LEVEL=
LOG_FORMAT=

DB_URI=
DB_USER=
DB_PASS=

JWT_KEY=
```

Run `make infra-down` or `make cleanup` in the end.

### Run the service in Docker

Run `make webservice` and `make cleanup` in the end.

### Clean code steps

`make lint`

`make run-test`

`make run-integration`

## Endpoints

### User related

<details>
    <summary>Endpoints</summary>

#### Sign Up

* __URI:__ _/sign-up_
* __Method:__ _POST_
* __Headers:__ - <br/>


* __URL params:__ - <br/>
* __Query params:__ - <br/>
* __Data params:__  <br/>
    ``` 
    {
       "username": "username",
       "first_name": "First",
       "last_name": "Last",
       "email": "email@abc.com",
       "password": "password"
    }
     ```


* __Responses:__
  * __Success response:__
    * Code: 201 Created <br/>
    * Content: -
  * __Failure response:__
    * Code: 400 Bad Request <br/>
    * Content: `user already exists`

#### Login

* __URI:__ _/login_
* __Method:__ _POST_
* __Headers:__ <br/>
  * `X-Request-Id` - correlation id as uuid
  * `Authorization` - Basic auth


* __URL params:__ - <br/>
* __Query params:__ - <br/>
* __Data params:__ - <br/>


* __Responses:__
  * __Success response:__
    * Code: 200 OK <br/>
    * Content: -
    * Sets cookie.
  * __Failure response:__
    * Code: 400 Bad Request <br/>
    * Content: `bad request`
    * If `Authorization` Header is missing.
  * __Failure response:__
    * Code: 401 Unauthorized <br/>
    * Content: `invalid username or password`

#### Logout

* __URI:__ _/logout_
* __Method:__ _POST_
* __Headers:__ <br/>
  * `X-Request-Id` - correlation id as uuid


* __URL params:__ - <br/>
* __Query params:__ - <br/>
* __Data params:__ - <br/>


* __Responses:__
  * __Success response:__
    * Code: 200 OK <br/>
    * Content: -
    * Removes cookie.

#### Refresh

* __URI:__ _/auth/refresh_
* __Method:__ _POST_
* __Headers:__ <br/>
  * `X-Request-Id` - correlation id as uuid
  * The cookie is required.

* __URL params:__ - <br/>
* __Query params:__ - <br/>
* __Data params:__ - <br/>


* __Responses:__
  * __Success response:__
    * Code: 200 OK <br/>
    * Content: -
    * Updates cookie.
  * __Failure response:__
    * Code: 401 Unauthorized <br/>
    * Content: `unauthorized`
    * If the token is missing.

</details>
