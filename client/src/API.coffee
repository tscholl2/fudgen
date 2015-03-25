settings = require './Settings.coffee'
window.$ = require 'jquery'

module.exports =
	getRecipe: (o) ->
		callback = o.callback
		$.ajax
			type: 'GET'
			url:  "#{settings.url}/api"
			xhrFields:
				withCredentials: false #use true if server has whitelist
			success: (response) ->
				console.log "got recipe:"
				console.log response
				callback? response
			dataType: 'json'
