Backbone = require 'backbone'
Backbone.$ = require 'jquery'
Backbone._ = require 'underscore'

sample_recipe = '
{
    "error":"",
    "recipe":
    {
        "steps":[
            {"name":"grilled cheese","op":"grill","id":0,"time":{"u":"minute","a":10},"reqs":[1,2],"notes":"flip at some point"},
            {"name":"butter","id":1,"data":{"Com_Desc":"Unsalted Butter","Long_Desc":"Margarine-like, butter-margarine blend, 80% fat, stick, without salt","ManufacName":"","NDB_No":"42307","Shrt_Desc":"MARGARINE-LIKE,BUTTER-MARGARINE BLEND,80% FAT,STK,WO/ SALT","price":"0.095234"},"quant":{"u":"tablespoon","a":0.25},"notes":""},
            {"name":"sandwich","op":"layer","id":2,"time":{"u":"minute","a":2},"reqs":[3,4],"notes":""},
            {"name":"bread","id":3,"data":{"Com_Desc":"Breadfruit","Long_Desc":"Breadfruit, raw","ManufacName":"","NDB_No":"09059","Shrt_Desc":"BREADFRUIT,RAW","price":"0.056504"},"quant":{"u":"fruit, small","a":0.5},"notes":"use whole wheat"},
            {"name":"cheese","id":4,"data":{"Com_Desc":"Oaxaca cheese","Long_Desc":"Cheese, mexican, queso asadero","ManufacName":"","NDB_No":"01166","Shrt_Desc":"CHEESE,MEXICAN,QUESO ASADERO","price":"0.055248"},"quant":{"u":"cubic inch","a":3},"notes":""}
        ],
        "title":"untited",
        "nutr":{"10:0":{"u":"g","a":0.30752855},"12:0":{"u":"g","a":0.28564591},"14:0":{"u":"g","a":1.706428349},"16:0":{"u":"g","a":5.264841404},"16:1 undifferentiated":{"u":"g","a":0.5095547939999999},"18:0":{"u":"g","a":2.225949786},"18:1 undifferentiated":{"u":"g","a":5.010045224},"18:2 undifferentiated":{"u":"g","a":0.93634537},"18:3 undifferentiated":{"u":"g","a":0.22792484200000002},"18:4":{"u":"g","a":0},"20:1":{"u":"g","a":0},"20:4 undifferentiated":{"u":"g","a":0},"20:5 n-3 (EPA)":{"u":"g","a":0},"22:1 undifferentiated":{"u":"g","a":0},"22:5 n-3 (DPA)":{"u":"g","a":0},"22:6 n-3 (DHA)":{"u":"g","a":0},"4:0":{"u":"g","a":0.492826214},"6:0":{"u":"g","a":0.259826533},"8:0":{"u":"g","a":0.14086043},"Alanine":{"u":"g","a":0.27054},"Alcohol, ethyl":{"u":"g","a":0},"Arginine":{"u":"g","a":0.4104},"Ash":{"u":"g","a":2.7397340000000003},"Aspartic acid":{"u":"g","a":0.7965000000000001},"Caffeine":{"u":"mg","a":0},"Calcium, Ca":{"u":"mg","a":366.135076},"Carbohydrate, by difference":{"u":"g","a":14.5895802},"Carotene, alpha":{"u":"µg","a":0},"Carotene, beta":{"u":"µg","a":24.709870000000002},"Cholesterol":{"u":"mg","a":59.953096},"Choline, total":{"u":"mg","a":13.574505000000002},"Copper, Cu":{"u":"mg","a":0.054544835},"Cryptoxanthin, beta":{"u":"µg","a":0},"Cystine":{"u":"g","a":0.05993999999999999},"Energy":{"u":"kJ","a":1122.5288679999999},"Fatty acids, total monounsaturated":{"u":"g","a":5.533573544000001},"Fatty acids, total polyunsaturated":{"u":"g","a":1.164810212},"Fatty acids, total saturated":{"u":"g","a":10.704660168},"Fiber, total dietary":{"u":"g","a":2.3520000000000003},"Folate, DFE":{"u":"µg","a":11.113934},"Folate, food":{"u":"µg","a":11.113934},"Folate, total":{"u":"µg","a":11.113934},"Folic acid":{"u":"µg","a":0},"Glutamic acid":{"u":"g","a":2.73888},"Glycine":{"u":"g","a":0.13283999999999999},"Gm_Wgt":{"u":"gram","a":105.69669999999999},"Histidine":{"u":"g","a":0.37476},"Iron, Fe":{"u":"mg","a":0.5379270300000001},"Isoleucine":{"u":"g","a":0.6922200000000001},"Leucine":{"u":"g","a":1.17762},"Lutein + zeaxanthin":{"u":"µg","a":10.56},"Lycopene":{"u":"µg","a":0},"Lysine":{"u":"g","a":0.85368},"Magnesium, Mg":{"u":"mg","a":26.113934},"Manganese, Mn":{"u":"mg","a":0.04824},"Methionine":{"u":"g","a":0.32502000000000003},"Niacin":{"u":"mg","a":0.5308490100000001},"Pantothenic acid":{"u":"mg","a":0.3441},"Phenylalanine":{"u":"g","a":0.65508},"Phosphorus, P":{"u":"mg","a":254.470241},"Potassium, K":{"u":"mg","a":283.007779},"Proline":{"u":"g","a":1.29762},"Protein":{"u":"g","a":12.7508703},"Retinol":{"u":"µg","a":58.090655999999996},"Riboflavin":{"u":"mg","a":0.13629868},"Selenium, Se":{"u":"µg","a":8.1364835},"Serine":{"u":"g","a":0.63126},"Sodium, Na":{"u":"mg","a":382.695076},"Sugars, total":{"u":"g","a":6.8298000000000005},"Theobromine":{"u":"mg","a":0},"Thiamin":{"u":"mg","a":0.06450967},"Threonine":{"u":"g","a":0.42294},"Total lipid (fat)":{"u":"g","a":18.354036900000004},"Tryptophan":{"u":"g","a":0.14526000000000003},"Tyrosine":{"u":"g","a":0.58206},"Valine":{"u":"g","a":0.7926},"Vitamin A, IU":{"u":"IU","a":234.830959},"Vitamin A, RAE":{"u":"µg","a":59.975972999999996},"Vitamin B-12":{"u":"µg","a":0.5436967},"Vitamin B-12, added":{"u":"µg","a":0},"Vitamin B-6":{"u":"mg","a":0.07698967000000001},"Vitamin C, total ascorbic acid":{"u":"mg","a":13.9236967},"Vitamin D":{"u":"IU","a":11.783604},"Vitamin D (D2 + D3)":{"u":"µg","a":0.2810901},"Vitamin D3 (cholecalciferol)":{"u":"µg","a":0.27},"Vitamin E (alpha-tocopherol)":{"u":"mg","a":0.362435},"Vitamin E, added":{"u":"mg","a":0},"Vitamin K (phylloquinone)":{"u":"µg","a":4.1569603},"Water":{"u":"g","a":57.2624786},"Zinc, Zn":{"u":"mg","a":1.68913934},"servings":{"u":"serving","a":3.76405}},
        "price":0.20698599999999998
    },
    "sched":[[[2,120],[0,600]]]
}
'
window.r = sample_recipe


class RecipeModel extends Backbone.Model
    defaults:
        error: ""
        schedule: []
        steps: []
        title: ""
        nutrition: {}
    initialize: ->
        console.log "recipe model created"
    parse: (message) ->
        data = JSON.parse message
        if (not data?) or (data.error? and data.error != "")
            @set "error", data.error
        else
            @set "title", data.recipe.title
            @set "steps", data.recipe.steps
            @set "nutrition", data.recipe.nutr
            @set "price", data.recipe.price
            @set "schedule", data.sched
recipe_model = new RecipeModel()

class RecipeView extends Backbone.View
    model: recipe_model
    el: null
    initialize: ->
        @listenTo @model, "change", @resizeSVG
    resizeSVG: ->
        Backbone.$ "#main-canvas"
        .height Backbone.$("#main-screen")[0].scrollHeight
recipe_view = new RecipeView()

class ScheduleView extends Backbone.View
    model: recipe_model
    el: $ "#schedule_container"
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
                    output += "<li class='list-group-item'>Do step #{job[0]} for #{parseInt job[1]/60} minutes.</li>"
            output += "</ul></div>"
        output += "</div>"
        @$el.html output
schedule_view = new ScheduleView()

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
    model: recipe_model
    el: $ "#steps_container"
    initialize: ->
        @listenTo @model, "change:steps", @render
    isIngrediant: (o) ->
        return not o.reqs?
    render: ->
        output = "<hr style='border-top:1px solid #666;'><h2>Steps</h2><div class='row'>"
        step_data = @model.get "steps"
        output += "<ul class='list-group'>"
        for s in step_data
            output += "<li class='list-group-item'>"
            if @isIngrediant s
                output += "[#{s.id}] - (#{s.name}): #{s.quant.a} #{s.quant.u} of #{s.data.Long_Desc}."
            else
                output += "[#{s.id}] - (#{s.name}): #{s.op} steps #{s.reqs} for #{s.time.a} #{s.time.u}s."
                if s.notes? and s.notes != ""
                    output += " Note: #{s.notes}"
            output += "</li>"
        output += "</ul>"
        @$el.html output
steps_view = new StepsView()


class TitleView extends Backbone.View
    model: recipe_model
    el: $ "#title_container"
    initialize: ->
        @listenTo @model, "change:title", @render
    render: ->
        output = "<h2>#{@model.get 'title'}</h2>"
        @$el.html output
title_view = new TitleView()


class PriceView extends Backbone.View
    model: recipe_model
    el: $ "#price_container"
    initialize: ->
        @listenTo @model, "change:price", @render
    render: ->
        output = "<hr style='border-top:1px solid #666;'><h2>Price: $#{@model.get 'price'}</h2>"
        @$el.html output
price_view = new PriceView()


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
    model: recipe_model
    el: $ "#nutrition_container"
    initialize: ->
        @listenTo @model, "change:nutrition", @render
    render: ->
        output = "<hr style='border-top:1px solid #666;'><h2>Nutrition</h2><div class='row'>"
        nutr_data = @model.get "nutrition"
        output += "<ul class='list-group align-left'>"
        for k,v of nutr_data
            output += "<li class='list-group-item'>#{v.a} #{v.u} of #{k}</li>"
        output += "</ul>"
        @$el.html output
nutr_view = new NutrView()



window.m = recipe_model
