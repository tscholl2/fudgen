sqlite3 = require 'sqlite3'

db = new sqlite3.Database '../data/db'

randomFood = (callback) ->
	output = {}
	db.get "select * from FOOD_DES where rowid = (abs(random()) % (select max(rowid)+1 from FOOD_DES))", (e,row) ->
	#db.get "select * from FOOD_DES where NDB_No='02001'", (e,row) -> #for testing
		if e?
			return callback? e
		output.FOOD_DES = row
		db.all "select * from FOOTNOTE where NDB_No=?", output.FOOD_DES.NDB_No, (e,rows) ->
			if e?
				return callback? e
			output.FOOTNOTE = rows
			db.all "select * from NUT_DATA where NDB_No=?", output.FOOD_DES.NDB_No, (e,rows) ->
				if e?
					return callback? e
				output.NUT_DATA = rows
				db.all "select * from WEIGHT where NDB_No=?",output.FOOD_DES.NDB_No, (e,rows) ->
					if e?
						return callback? e
					output.WEIGHT = rows
					db.all "select * from FD_GROUP where FdGrp_Cd=?",output.FOOD_DES.FdGrp_Cd, (e,rows) ->
						if e?
							return callback? e
						output.FD_GROUP = rows
						db.all "select * from LANGUAL where NDB_No=?",output.FOOD_DES.NDB_No, (e,rows) ->
							if e?
								return callback? e
							output.LANGUAL = rows
							db.all "select * from LANGDESC where Factor_Code in (#{('\''+r.Factor_Code+'\'' for r in output.LANGUAL).join(',')})", (e,rows) ->
								if e?
									return callback? e
								output.LANGDESC = rows
								db.all "select * from NUTR_DEF where Nutr_No in (#{('\''+r.Nutr_No+'\'' for r in output.NUT_DATA).join(',')})", (e,rows) ->
									if e?
										return callback? e
									output.NUTR_DEF = rows
									return callback? null, output


randomFood (e,food) ->
	console.log food
