# For fresh run
Execute in terminal
```
> make fresh_run 
> make build-dbs
```
And wait till databases get ready for connections\
Then in separate terminal
```
> make db_init
```
Stop databases containers and then start all containers 
```
> docker-compose down
> docker-compose up --build -d
```

Open url in browser http://localhost:8080/

You are all set:)
