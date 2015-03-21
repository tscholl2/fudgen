console.log "Hello"

settings = require './settings.coffee'
console.log settings

SM = require './ScreenManager.coffee'
console.log SM

B = require './Bunny.coffee'

M = require './Models.coffee'

window.b = new B(document.getElementById("main-canvas"))

#if this is running then everything has been loaded

_fn = ->
    b.gp.center window.innerWidth/2, window.innerHeight - 17
    SM.Switch "main-screen"
    b.ActNatural true

setTimeout _fn, 750
