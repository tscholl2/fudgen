console.log "Hello"

settings = require './settings.coffee'
console.log settings

SM = require './ScreenManager.coffee'
console.log SM

B = require './Bunny.coffee'

window.b = new B(document.getElementById("main-canvas"))

console.log

#if this is running then everything has been loaded



_fn = ->
    b.gp.center window.innerWidth/2, window.innerHeight/2
    SM.Switch "main-screen"
setTimeout _fn, 750
