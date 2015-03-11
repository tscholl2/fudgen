window.$ = require('jquery')
console.log window.$


class ScreenManager

	constructor: ->
		@name = "loading-screen"

	Switch: (name,callback) ->
		$("##{@.name}").fadeOut =>
			@.name = name
			$("##{@.name}").fadeIn callback

	Current: ->
		return @name




module.exports = new ScreenManager()
