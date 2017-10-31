[![Build Status](https://travis-ci.org/w32blaster/fyzon.svg?branch=master)](https://travis-ci.org/w32blaster/fyzon)

# Fyzon - Translation tool

Lightweight. Simple. Free.

# Fyzon resource
If you use [ConcourseCI](http://concourse.ci/) to build your projects, then consider to use my [Fyzon Resource](https://github.com/w32blaster/fyzon-resource) to automate translations downloading.

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


# Few screenshots here:

Few screenshots for you.

## Main page
![The Main page](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_069.jpg)


## Project page
![Selected project page](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_070.jpg)


## Upload new .properties file to import
![Import new file popup](https://raw.githubusercontent.com/w32blaster/monsieur-traducteur/master/docs/Selection_072.jpg)


# Run
Do you want to try it? Huh, that's super easy with our [docker image](https://hub.docker.com/r/w32blaster/fyzon/). 
Simply run:

```
docker run -p 8080:8080 w32blaster/fyzon
```

and navigate to http://localhost:8080. Login with [user](https://github.com/w32blaster/fyzon/blob/master/models.user.go#L21) *user1* and password *pass1* (registration page is ongoing). Enjoy!

Or, if you want to persist the database outside of container (which is recommended for production), then specify volume when you start the container:

```
docker run -p 8080:8080 --volume /your/directory:/go/bin/db w32blaster/fyzon
```

here is the same configuration, but for docker-compose:

```
 fyzon:
   container_name: fyzon
   restart: always
   image: w32blaster/fyzon:latest
   volumes:
     - /your/directory:/go/bin/db
   ports:
     - 8080:8080

```

and then you need to set up the database (only the first time). The schema is stored in the container, so you need only call the command:

```

docker exec -it fyzon /go/bin/importDb.sh

```

After that you can restart the container and DB will always stay in the `/your/directory` directory.

# How to download generated file?

You can download the file with translations using API. Say, you have a CI/CD pipeline and before you build your application you might want to 
dowload fresh translations. Use old plain CURL to do that:

## JSON format

JSON format could be used in Go projects with [go-i18n](https://github.com/nicksnyder/go-i18n) library. 

```
curl http://localhost:8080/api/project/3/file/gb/json > en.all.json
```

where:
  * "3" is the ID of project
  * "gb" is desired language (country) to be exported
  * "json" is type of file

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

where "%3D" is encoded "=" symbol. The result will be something like that:

```
button.submit= Submit
button.cancel= Cancel
...
```

----

# How to start the development on your computer

Wow, you are reading this paragraph, that means you decided to bless this badly written code with your attention! :grin: Welcome, mate!


To start coding, you need only *sqlite3*, *go*, *npm* and *govendor* installed on your computer. Clone the code and navigate to the project folder, then:

1) import the database schema

```
sqlite3 db/trans.sqlite3 < db/schema.sql
```

2) download and install all the golang dependencies (assuming that you have [govendor](https://github.com/kardianos/govendor) installed)

```
govendor fetch -v +out

```

3) rebuild [SemanticUI](https://semantic-ui.com/) theme

```
npm install semantic-ui --save
cd semantic
gulp build
```

And here you go! Just compile with command `go build` and run it with `./fyzon`! Or, in dev mode execute the command `go run !(*_test*).go`. 
Please, enjoy the development, I would be happy to accept an PR from you! 

Good luck!
