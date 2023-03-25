# Luzifer / mysqlapi

This repo contains a simple-ish web-application to translate HTTP POST requests into MySQL queries.

## Why?!?

I had the requirement to do SQL queries from fairly simple scripts without the possibility to add a MySQL client. As HTTP calls are possible in nearly every environement the idea was to have an API to execute arbitrary SQL statements over a JSON POST-API.

## Security

**⚠⚠⚠ NEVER EVER LEAVE THIS OPEN TO THE INTERNET! ⚠⚠⚠**

Having stated that as clearly as possible: This API does not limit the type of queries being executed. The only thing saving you might be the permissions of the user you configured in the DSN given to the tool. If you gave it global admin permissions, well - you've just handed over your database server.

In general make sure you understood what is possible using this and limit access to an absolute minimum. Your data got lost / leaked? I did warn you.

## How to use?

```
POST /{database}
Content-Type: application/json

[...]
```

```console
# mysqlapi --help
Usage of mysqlapi:
      --dsn string         MySQL DSN to connect to: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
      --listen string      Port/IP to listen on (default ":3000")
      --log-level string   Log level (debug, info, warn, error, fatal) (default "info")
      --version            Prints current version and exits

# mysqlapi \
    --dsn "limiteduser:verysecretpass@tcp(mydatabase.cluster.local:3306)/?charset=utf8mb4&parseTime=True&loc=Local" \
    --listen 127.0.0.1:7895
INFO[0000] mysqlapi started                              addr="127.0.0.1:7895" version=dev

# curl -s --data-binary @select.json localhost:7895/mysqlapi_test | jq .
```

**Request format**

```json
[
  ["SELECT * FROM testtable"],
  ["INSERT INTO testtable (name, age, birthday) VALUES (?, ?, ?)", "Karl", 45, "1999-02-05T02:00:00"],
  ["SELECT * FROM testtable WHERE name = ?", "Karl"],
  ["DELETE FROM testtable WHERE name = ?", "Karl"]
]
```

**Response format**

```json
[
  null,
  null,
  [
    {
      "age": 45,
      "birthday": "1999-02-05T02:00:00+01:00",
      "id": 1,
      "name": "Karl"
    }
  ],
  null
]
```
