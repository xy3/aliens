# Alien Invasion

**Author: Theodore Coyne Morgan | March 2022**

This program simulates an alien invasion on a map of cities.

## Usage

Running the simulation:
```shell
git clone git@github.com:xy3/aliens.git
cd aliens
go run cmd/main.go -h
```

When you run the program it will give you a summary of program arguments, flags and a brief description of the application.

For example:
```
Simulate an alien invasion using a configurable amount of aliens and a custom map.

You may use the -h flag to view help for this program. To modify the default config,
modify 'config.json' in the program's working directory.
Author:
Theodore Coyne Morgan | March 2022

Usage:
  aliens [alien_count] [flags]

Flags:
  -n, --alien-names string   text file with names of simulation aliens on each line (default "alien-names.txt")
  -d, --debug                debug mode to view more logging
  -h, --help                 help for aliens
  -m, --map string           text file with cities and routes for the simulation (default "map.txt")
```

An example use of this program is as follows:

```shell
go run cmd/main.go 10 -n alien-names.txt -m map.txt -d
# run the application using 10 random aliens
# -n: use the alien-names.txt file to choose names from
# -m: use map.txt to load in the world map
# -d: view debug logs
```

The contents of the map file could be as follows:

```
Paris east=London south=Rome
Madrid north=London south=Milan east=Tokyo west=Rome
Tokyo north=Dublin west=Madrid east=Lisbon
Dubai west=Dublin south=Lisbon
Milan west=Amsterdam east=Naples south=Berlin
Budapest north=Amsterdam east=Berlin
```

Each line defines a city with connections in each compass direction separated by a space.

## Configuration

If you wish to configure the application using a json file, you may do so by placing a file called `config.json` in the program's working directory.

It will also automatically be created after running the program once.

Here is an example value for the config:

```json
{
    "MaxAlienMoves": 100000,
    "MapFile": "map.txt",
    "AlienNamesFile": "alien-names.txt",
    "DebugMode": false
}
```

The names should be self-explanatory. 

## Running tests

From the repo root run:

```shell
go test ./... -cover
```

## Example output

```go
INFO[0000] Loading config from: /Users/tcoynemorgan/go/src/github.com/xy3/aliens/config.json 
INFO[0000] PARSED SIMULATION WORLD MAP                  
Amsterdam east=Milan south=Budapest
Berlin north=Milan west=Budapest
Budapest north=Amsterdam east=Berlin
Dubai south=Lisbon west=Dublin
Dublin east=Dubai south=Tokyo
Lisbon north=Dubai west=Tokyo
London south=Madrid west=Paris
Madrid north=London east=Tokyo south=Milan west=Rome
Milan east=Naples south=Berlin west=Amsterdam
Naples west=Milan
Paris east=London south=Rome
Rome north=Paris east=Madrid
Tokyo north=Dublin east=Lisbon west=Madrid
INFO[0000] CREATED 8 RANDOM ALIENS SUCCESSFULLY         
INFO[0000] Deployed Galen to Rome                       
INFO[0000] Deployed Korben to Amsterdam                 
INFO[0000] Deployed Soval to Budapest                   
INFO[0000] Deployed Vaastot to Dublin                   
INFO[0000] Deployed Malcom to London                    
INFO[0000] Deployed Jaina to Madrid                     
INFO[0000] Deployed Tasha to Berlin                     
INFO[0000] Deployed Tharen to Dubai                     
INFO[0000] Madrid has been destroyed by Galen and Jaina  destroyedCity=Madrid opponents="Galen vs Jaina"
INFO[0000] Milan has been destroyed by Tasha and Korben  destroyedCity=Milan opponents="Tasha vs Korben"
INFO[0000] Tokyo has been destroyed by Vaastot and Tharen  destroyedCity=Tokyo opponents="Vaastot vs Tharen"
INFO[0000] 10000 DAYS HAVE PASSED                       
INFO[0000] 10000 DAYS HAVE PASSED                       
INFO[0001] 10000 DAYS HAVE PASSED                       
INFO[0001] 10000 DAYS HAVE PASSED                       
INFO[0001] 10000 DAYS HAVE PASSED                       
INFO[0002] 10000 DAYS HAVE PASSED                       
INFO[0002] 10000 DAYS HAVE PASSED                       
INFO[0002] 10000 DAYS HAVE PASSED                       
INFO[0003] 10000 DAYS HAVE PASSED                       
INFO[0003] 10000 DAYS HAVE PASSED                       
INFO[0003] ==== SIMULATION RESULTS: ====                
INFO[0003] Result                                        daysPassed=107291 dead=6 exhausted=2 stuck=0
INFO[0003] The world map that still remains is:         
Amsterdam east=Milan south=Budapest
Berlin north=Milan west=Budapest
Budapest north=Amsterdam east=Berlin
Dubai south=Lisbon west=Dublin
Dublin east=Dubai south=Tokyo
Lisbon north=Dubai west=Tokyo
London south=Madrid west=Paris
Naples west=Milan
Paris east=London south=Rome
Rome north=Paris east=Madrid
```

## Possible design changes / future work

There are a number of various changes that could be made to make this task more exciting or efficient:

- Make each Alien run on its own go routine and handle interactions with the world map using locks
- Make the world map exist somewhere on the network and communicate with it remotely
- Display the world map via the console as a spacial representation of the input map
- Add results analysis
- Preemptively shutdown the simulation if Aliens are stuck in a loop

## Challenges

- Managing mutability with the world map and pointers
- Writing unit tests for the simulation

---

### Thanks for this task!

Email: hi@theodore.ie
