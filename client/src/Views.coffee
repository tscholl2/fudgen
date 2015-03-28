window.jQuery = require 'jquery'
window.$ = jQuery
require 'bootstrap/dist/js/bootstrap.js'
Backbone = require 'backbone'
Backbone.$ = jQuery
Backbone._ = require 'underscore'



class RecipeView extends Backbone.View
    initialize: ->
        @listenTo @model, "change", @resizeSVG
    resizeSVG: ->
        Backbone.$ "#main-canvas"
        .height Backbone.$("#main-screen")[0].scrollHeight

class ScheduleView extends Backbone.View
    initialize: ->
        @listenTo @model, "change:schedule", @render
    getPerson: (i) ->
        p = "ABCDEFGHIJKLMNOPQRSTUVWXYZ".split ""
        return "#{p[i % p.length]}#{if i < p.length then '' else parseInt i / p.length}"
    render: ->
        output = "<hr style='border-top:1px solid #666;'><h2>Schedule</h2><div class='row'>"
        schedule_data = @model.get "schedule"
        for arr,i in schedule_data
            output += "<div class='col-lg-3 col-md-6 col-sm-12'><h4>Person #{@getPerson i}</h4><ul class='list-group'>"
            for job in arr
                if job[0] == -1
                    output += "<li class='list-group-item'>Hang out for #{parseInt job[1]/60} minutes.</li>"
                else
                    output += "<li class='list-group-item'><b>Step #{job[0]}</b> for #{parseInt job[1]/60} minutes.</li>"
            output += "</ul></div>"
        output += "</div>"
        @$el.html output

# {
#     "name":"bread",
#     "id":3,
#     "data":{
#         "Com_Desc":"Breadfruit",
#         "Long_Desc":"Breadfruit, raw",
#         "ManufacName":"",
#         "NDB_No":"09059",
#         "Shrt_Desc":"BREADFRUIT,RAW",
#         "price":"0.056504"
#     },
#     "quant":{"u":"fruit, small","a":0.5},
#     "notes":"use whole wheat"
# },

# {
#     "name":"sandwich",
#     "op":"layer",
#     "id":2,
#     "time":{"u":"minute","a":2},
#     "reqs":[3,4],
#     "notes":""
# },

class StepsView extends Backbone.View
    initialize: ->
        @listenTo @model, "change:steps", @render
    isIngrediant: (o) ->
        return not o.reqs?
    render: ->

        hr = document.createElement "hr"
        $ hr
        .css "border-top", "1px solid #666"
        h2 = document.createElement "h2"
        h2.innerHtml = "Steps"

        step_data = @model.get "steps"
        table = document.createElement "table"
        $ table
        .addClass "table table-striped table-hover"
        .attr "id", "steps_table"
        table.innerHtml += "<tr><th>#</th><th>Name</th><th>Description</th></tr>"
        for s in step_data
            tr = document.createElement "tr"
            if @isIngrediant s
                $ tr
                .html "<th>#{s.id}</th><td>#{s.name}</td><td>#{s.quant.a} #{s.quant.u} of #{s.data.Long_Desc}</td>"
                .data "s", s
                .click ->
                    wrap = (label,value) ->
                        return "<div class='form-group'><label class='col-sm-2 control-label'>#{label}</label><div class='sm-col-10'><p class='form-control-static'>#{value}</p></div></div>"
                    s = $(@).data "s"
                    $ "#myModalLabel"
                    .html "Step #{s.id}"
                    body = "<form class='form-horizontal'>"
                    body += wrap "ID", s.id
                    body += wrap "Name", s.name
                    if s.notes != ""
                        body += wrap "Notes", s.notes
                    body += wrap "Quantity", "#{s.quant.a} #{s.quant.u}"
                    for k,v of s.data
                        if v !=  ""
                            body += wrap k, v
                    body += "</form>"
                    $ "#myModalBody"
                    .html body
                    $ "#myModal"
                    .modal()
            else
                $ tr
                .html "<th>#{s.id}</th><td>#{s.name}</td><td>#{s.op} steps <b>#{s.reqs}</b> for #{s.time.a} #{s.time.u}s</td>"
                .data "s", s
                .click ->
                    s = $(@).data "s"
                    $ "#myModalLabel"
                    .html "Step #{s.id}"
                    body = "ID = #{s.id}<br>"
                    body += "Name = #{s.name}<br>"
                    body += "Operation = #{s.op}"
                    body += "Requirements = #{s.reqs}"
                    body += "Notes = #{s.notes}<br>"
                    body += "Time = #{s.time.a} #{s.time.u}s<br>"
                    $ "#myModalBody"
                    .html body
                    $ "#myModal"
                    .modal()
            $ table
            .append tr
        $ @$el
        .html ""
        .append hr, h2, table


class TitleView extends Backbone.View
    initialize: ->
        @listenTo @model, "change:title", @render
    render: ->
        output = "<h2>#{@model.get 'title'}</h2>"
        @$el.html output


class PriceView extends Backbone.View
    initialize: ->
        @listenTo @model, "change:price", @render
    render: ->
        output = "<hr style='border-top:1px solid #666;'><h2>Price <font color='green'>$#{Math.round(@model.get('price') * 100) / 100}</font></h2>"
        @$el.html output


# "nutr":{
#     "Alanine":{"u":"g","a":0.27054},
#     "Alcohol, ethyl":{"u":"g","a":0},
#     "Arginine":{"u":"g","a":0.4104},
#     "Ash":{"u":"g","a":2.7397340000000003},
#     "Aspartic acid":{"u":"g","a":0.7965000000000001},
#     "Caffeine":{"u":"mg","a":0},
#     "Calcium, Ca":{"u":"mg","a":366.135076},
#     "Carbohydrate, by difference":{"u":"g","a":14.5895802},
#     "Carotene, alpha":{"u":"µg","a":0},
#     "Carotene, beta":{"u":"µg","a":24.709870000000002},
#     "Cholesterol":{"u":"mg","a":59.953096},
#     "Choline, total":{"u":"mg","a":13.574505000000002},
#     "Copper, Cu":{"u":"mg","a":0.054544835},
#     "Cryptoxanthin, beta":{"u":"µg","a":0},
#     "Cystine":{"u":"g","a":0.05993999999999999},
#     "Energy":{"u":"kJ","a":1122.5288679999999},
#     "Fatty acids, total monounsaturated":{"u":"g","a":5.533573544000001},
#     "Fatty acids, total polyunsaturated":{"u":"g","a":1.164810212},
#     "Fatty acids, total saturated":{"u":"g","a":10.704660168},
#     "Fiber, total dietary":{"u":"g","a":2.3520000000000003},
#     "Folate, DFE":{"u":"µg","a":11.113934},"Folate, food":{"u":"µg","a":11.113934},
#     "Folate, total":{"u":"µg","a":11.113934},"Folic acid":{"u":"µg","a":0},
#     "Glutamic acid":{"u":"g","a":2.73888},"Glycine":{"u":"g","a":0.13283999999999999},
#     "Gm_Wgt":{"u":"gram","a":105.69669999999999},"Histidine":{"u":"g","a":0.37476},
#     "Iron, Fe":{"u":"mg","a":0.5379270300000001},"Isoleucine":{"u":"g","a":0.6922200000000001},
#     "Leucine":{"u":"g","a":1.17762},"Lutein + zeaxanthin":{"u":"µg","a":10.56},
#     "Lycopene":{"u":"µg","a":0},"Lysine":{"u":"g","a":0.85368},
#     "Magnesium, Mg":{"u":"mg","a":26.113934},"Manganese, Mn":{"u":"mg","a":0.04824},
#     "Methionine":{"u":"g","a":0.32502000000000003},"Niacin":{"u":"mg","a":0.5308490100000001},
#     "Pantothenic acid":{"u":"mg","a":0.3441},"Phenylalanine":{"u":"g","a":0.65508},
#     "Phosphorus, P":{"u":"mg","a":254.470241},"Potassium, K":{"u":"mg","a":283.007779},
#     "Proline":{"u":"g","a":1.29762},"Protein":{"u":"g","a":12.7508703},
#     "Retinol":{"u":"µg","a":58.090655999999996},
#     "Riboflavin":{"u":"mg","a":0.13629868},
#     "Selenium, Se":{"u":"µg","a":8.1364835},
#     "Serine":{"u":"g","a":0.63126},
#     "Sodium, Na":{"u":"mg","a":382.695076},
#     "Sugars, total":{"u":"g","a":6.8298000000000005},
#     "Theobromine":{"u":"mg","a":0},
#     "Thiamin":{"u":"mg","a":0.06450967},
#     "Threonine":{"u":"g","a":0.42294},
#     "Total lipid (fat)":{"u":"g","a":18.354036900000004},
#     "Tryptophan":{"u":"g","a":0.14526000000000003},
#     "Tyrosine":{"u":"g","a":0.58206},
#     "Valine":{"u":"g","a":0.7926},
#     "Vitamin A, IU":{"u":"IU","a":234.830959},
#     "Vitamin A, RAE":{"u":"µg","a":59.975972999999996},
#     "Vitamin B-12":{"u":"µg","a":0.5436967},
#     "Vitamin B-12, added":{"u":"µg","a":0},
#     "Vitamin B-6":{"u":"mg","a":0.07698967000000001},
#     "Vitamin C, total ascorbic acid":{"u":"mg","a":13.9236967},
#     "Vitamin D":{"u":"IU","a":11.783604},
#     "Vitamin D (D2 + D3)":{"u":"µg","a":0.2810901},
#     "Vitamin D3 (cholecalciferol)":{"u":"µg","a":0.27},
#     "Vitamin E (alpha-tocopherol)":{"u":"mg","a":0.362435},
#     "Vitamin E, added":{"u":"mg","a":0},
#     "Vitamin K (phylloquinone)":{"u":"µg","a":4.1569603},
#     "Water":{"u":"g","a":57.2624786},
#     "Zinc, Zn":{"u":"mg","a":1.68913934},
#     "servings":{"u":"serving","a":3.76405}
# },

class NutrView extends Backbone.View
    initialize: ->
        @listenTo @model, "change:nutrition", @render
    parseUnit: (u) ->
        if u.charCodeAt(0) == 194
            return "&#181;g"
        if u == "gram"
            return "g"
        if u == "serving"
            return "Servings"
        return u
    render: ->
        output = "<hr style='border-top:1px solid #666;'><h2>Nutrition</h2>"
        nutr_data = @model.get "nutrition"
        output += "<div class='row'>"
        for k,v of nutr_data
            output += "<div class='col-xs-12 col-sm-6 col-md-4'><div class='panel panel-info'>"
            output +="<div class='panel-title' style='font-size:small;'>#{k}</div>"
            output += "<div class='panel-body'><h4>#{Math.round(v.a * 1000) / 1000} #{@parseUnit v.u}</h4></div>"
            output += "</div></div>"
        output += "</div>"
        @$el.html output

module.exports =
    Nutrition: NutrView
    Price: PriceView
    Schedule: ScheduleView
    Title: TitleView
    Steps: StepsView
    Recipe: RecipeView
