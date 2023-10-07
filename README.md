
Make sure docker is installed: https://docs.docker.com/desktop/

```
  docker-compose up
```
Did this pretty quickly with vanilla Go but decided to keep going with an interesting database.

Starts up a [surrealdb](https://surrealdb.com) database. My Go image connects to it.

Manually test with a browser or curl.  
```
http://localhost:3333/upload
```
Doing a GET on the upload path present an html form to upload a csv.

Otherwise with curl to upload a csv:
```
curl -X POST -F "csvfile=@seattle-weather.csv" http://localhost:3333/upload
```

to query a date in a browser:
```
http://localhost:3333/query?date=2012-02-06
```

curl
```
curl "http://localhost:3333/query?date=2012-02-06"
```

to query a weather type with a limit in a browser:
```
http://localhost:3333/query?weather=drizzle&limit=6
```

curl
```
curl "http://localhost:3333/query?weather=drizzle&limit=6"
```

