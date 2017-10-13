[![Build Status](https://travis-ci.org/w32blaster/fyzon.svg?branch=master)](https://travis-ci.org/w32blaster/fyzon)

# Fyzon - Translation tool

Lightweight. Simple. Free.

## Why
Ok, say, you have an application and you have to localyze it. Say, you have a `*.properies` file with some phrases in English
and you need to make similar files with another languages. You make them locally and your app works. Next, you want to
ask a translator to translate phrases and he/she struggles to connect to Git and edit `*.properies` documents, because he/she
is used to work only in MS Office. What to do?

## Alternatives?
Ok, there are number of professional translation services, that can help you. But they cost money and while your project grows, the price 
increases drammatically.

## Here comes Fyzon
Fyzon is simple as hell. You can install it on your server via only one command and send the URL link to your translator. 
Then you can download the latest copy of `.properties` file manually or using your CI via API. It is open source and always will be free.

# Install
Do you want to try it? Huh, that's super easy with our [docker image](https://hub.docker.com/r/w32blaster/fyzon/). 
Simply run:

```
docker run -p 8080:8080 w32blaster/fyzon
```

and navigate to http://localhost:8080. Login with [user](https://github.com/w32blaster/fyzon/blob/master/models.user.go#L21) *user1* and password *pass1* (registration page is ongoing). Enjoy!


# Few screenshots here:

Few screenshots for you.

## Main page
![The Main page](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_069.jpg)


## Project page
![Selected project page](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_070.jpg)


## Upload new .properties file to import
![Import new file popup](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_072.jpg)


# How to download generated file?

You can download the file with translations using API. Say, you have a CI/CD pipeline and before you build your application you might want to 
dowload fresh translations. Use old plain CURL to do that:

## JSON format

JSON format could be used in Go projects with [go-i18n](https://github.com/nicksnyder/go-i18n) library. 

```
curl http://localhost:8080/api/project/3/file/gb/json > en.all.json
```

where "3" is ID of project, "en" desired language to be exported and "json" is type of file

## Properties format

`.properties` format is very common in Java code, such as Spring.

```
curl http://localhost:8080/api/project/3/file/gb/properties > messages_en.properties
```

the result will be something like this:

```

button.submit: Submit
button.cancel: Cancel
...

```

or, here you can optionally specify delimeter you want to use (":" is default). In this example, we will user "=" as the delimeter

```
curl http://localhost:8080/api/project/3/file/gb/properties?delimeter=%3D > messages_en.properties
```

the result will be something like that:

```
button.submit=Submit
button.cancel=Cancel
...
```

----

# How to start the development on your computer

Wow, you read this paragraph, that means you decided to bless this badly written code with your attention! Welcome, mate!


To start coding, you need only *sqlite3*, *go* installed on your computer. Then,

1) and import the database schema

```
sqlite3 db/trans.sqlite3 < db/schema.sql
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

And here you go! Just compile with command `go build` and run it with `./fyzon`! Or, in dev mode `go run !(*_test*).go`. 
Please, join the development, I would be happy to accept an PR from you! 

Good luck!