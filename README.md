## Pokedex
Using the API available at https://pokeapi.co/, design a terminal application to retrieve data
about the Pokemon when queried.

The program must have the following functionality:

  ● Ability to cache the stored information locally (using a text file)

  ● Ability to search by Pokemon Name or ID.

  ● Display only the following information

    ○ Pokemon ID

    ○ Pokemon Name

    ○ Pokemon Type(s)

    ○ Pokemon Encounter Location(s) and method(s) in Kanto only

  ■ If there are no encounter location in Kanto, display ‘-’

    ○ Pokemon stats (speed, def, etc etc)

  ● If the stored information is over a week old, the data should be retrieved again from the API. If not, the data should be retrieved from the text file.

## Requirements
Go 1.16

## Run the app
```
make app
```
