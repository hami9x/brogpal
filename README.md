brogpal
=======

Companion project of the Wade framework - Brogrammer's palace

#Running
##Install Gopherjs:

    go get -u github.com/gopherjs/gopherjs
    go install github.com/gopherjs/gopherjs

##Install fresh:

    go get -u github.com/pilu/fresh
    go install github.com/pilu/fresh
    
You should have `gopherjs` and `fresh` in $PATH as commands now.  

##Install javascript dependencies:
- Install [bower](http://bower.io)
- Go to public/, run:

    bower install

##Actual run
Go to this project's directory, run:
    
    ./run_fresh

This one runs the server in server/, it's basically `cd server && ./fresh`. *Fresh* compiles the server and waits for any changes to the server go files and automatically recompile.
Make new terminal tab/window, run:

    ./run_gopherjs
This one actually runs something like `gopherjs -w <output directory...>` which runs gopherjs in watch mode, compiles the client go files to javascript, waiting for any changes to those go files and automatically recompile.

That's all for running. The site is typically served at http://localhost:3000.

#Development docs
[how wade works](https://github.com/phaikawl/wade)
