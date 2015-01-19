sqlite3 = require 'sqlite3'

db = new sqlite3.Database '../data/db'

randomFood = (callback) ->
	db.get "
	SELECT FOOD_DES.*,'^^',FOOTNOTE.*,'^^',NUT_DATA.*,'^^',WEIGHT.*,'^^',FD_GROUP.*,'^^',LANGUAL.*,'^^',LANGDESC.*,'^^',NUTR_DEF.* FROM FOOD_DES
	LEFT JOIN FOOTNOTE ON FOOD_DES.NDB_No=FOOTNOTE.NDB_No
	LEFT JOIN NUT_DATA ON FOOD_DES.NDB_No=NUT_DATA.NDB_No
	LEFT JOIN WEIGHT ON FOOD_DES.NDB_No=WEIGHT.NDB_No
	LEFT JOIN FD_GROUP ON FOOD_DES.FdGrp_Cd=FD_GROUP.FdGrp_Cd
	LEFT JOIN LANGUAL ON FOOD_DES.NDB_No=LANGUAL.NDB_No
	LEFT JOIN LANGDESC ON LANGUAL.Factor_Code=LANGDESC.Factor_Code
	LEFT JOIN NUTR_DEF ON NUT_DATA.Nutr_No=NUTR_DEF.Nutr_No
	WHERE FOOD_DES.rowid = (ABS(RANDOM()) % (SELECT MAX(FOOD_DES.rowid)+1 FROM FOOD_DES));
	", callback

randomFood (e,row) ->
	console.log e
	console.log row
