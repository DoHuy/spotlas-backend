# Spotlas Backend Developer Test

```azure
- task 1 (in folder task1 is *.sql file)
- task 2 (in folder task2)    
```
## How to run ?

```azure
 - step 1: install docker and docker-compose
 - step 2: cd spotlas && docker-compose up -d
 - step 3: dump Backend-Developer-Test/spots.sql to database
 - step 4: docker restart  spotlas_api_1    
```

### A. Task 1
You can run query in sql editor

### B. Task 2

You can test with api template below:

notes: type=1 (which is circle), type=2 (which is square)
```azure
curl -X GET \
  'http://localhost:3000/api/spots?latitude=52.21055029999999&longitude=0.1177139&radius=1000000&type=2' \
  -H 'cache-control: no-cache' \
  -H 'postman-token: 747af107-d489-f576-aafd-10092a32a77d'
```