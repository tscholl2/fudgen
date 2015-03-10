window.$ = require('jquery')
console.log window.$


class ScreenManager

	constructor: ->
		@name = "loading-screen"

	Switch: (name,callback) ->
		self = @
		$("##{self.name}").fadeOut ->
			self.name = name
			$("##{self.name}").fadeIn callback

	Current: ->
		return @name




module.exports = new ScreenManager()