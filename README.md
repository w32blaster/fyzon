# Translation tool

Lightweight. Simple. Free. 

Few screenshots here:

## Main page
![The Main page](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_069.jpg)


## Project page
![Selected project page](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_070.jpg)


## Upload new .properties file to import
![Import new file popup](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_072.jpg)

# How to start the development on your computer

Prerequirements: sqlite3, go

1) and import the database schema

```
sqlite3 trans.sqlite3 < db/schema.sql
```

2) download and install all the golang dependencies

```
go get -u -v github.com/gin-gonic/gin
go get -u -v github.com/mattn/go-sqlite3
go get -u -v github.com/stretchr/testify

```

3) rebuild SemanticUI theme

```
npm install semantic-ui --save
cd semantic
gulp build
```

