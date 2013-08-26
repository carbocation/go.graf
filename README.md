GRAF
====

`go get github.com/carbocation/go.graf`

GRAF Recursively Arranged Forum is a golang toolkit that contains all of the necessary parts for building a live-updating threaded forum 
using golang and a database (currently, only postgres is explicitly supported). This is the source code for http://askgolang.com .

To get started, take a look at the "example" directory. This example assumes you are using OS X for 
development and Postgres as your database server. If so, you can create what is needed in Postgres 
with "example/forum.sql", and then you can compile the project from that path with `./compile.sh` . 
(That will get you a dev.forum.osx to run locally and will also produce a prod.forum.linux, assuming that 
your actual server is a linux box. You can change this by modifying the contents of compile.sh .) 

If you just want to see a live-loading websocked-based chat functionality, check out [Go.Websocket-Chat](https://github.com/carbocation/go.websocket-chat), 
a library that is used in this project to provide comment live-loading.

The LICENSE specifies the terms; essentially, this project is licensed under an MIT-style license.

Please report any issues to me [here on Github](https://github.com/carbocation/go.graf).