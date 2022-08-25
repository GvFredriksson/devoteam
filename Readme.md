# Devoteam code assigment!

To run the service simply run:
```
make build
make up
```
If you want to see the terminal output instead run:
```
make build
make up_live
```

To run the test suit run:
```
make build
make test
```

If any problems with dependencies occur run:
```
make tidy
make verify
```


When the server is running, open a new terminal and curl request against localhost to test the service. Examples:

Ex1
```
curl -X POST -d 'SizeX=5&SizeY=5' http://localhost:8090/
curl -X POST -d 'X=0&Y=0&Direction=E' http://localhost:8090/initiate_robot
curl -X POST -d 'Moves=RFLFFLRF' http://localhost:8090/move_robot
```

Ex2
```
curl -X POST -d 'SizeX=5&SizeY=5' http://localhost:8090/
curl -X POST -d 'X=1&Y=2&Direction=N' http://localhost:8090/initiate_robot
curl -X POST -d 'Moves=RFRFFRFRF' http://localhost:8090/move_robot
```


Author GvFredriksson
