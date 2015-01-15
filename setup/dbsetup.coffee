fs = require 'fs'
sqlite3 = require 'sqlite3'

db = new sqlite3.Database '../data/db'

cleanField = (raw) ->
	if raw.length > 1 and raw[0] == '~' and raw[raw.length-1] == '~'
		return raw.substr 1, raw.length-2
	else
		return raw

#MAIN

main = () ->
	file_names = [
		'DATA_SRC'
		'DATSRCLN'
		'DERIV_CD'
		'FD_GROUP'
		'FOOD_DES'
		'FOOTNOTE'
		'LANGDESC'
		'LANGUAL'
		'NUT_DATA'
		'NUTR_DEF'
		'SRC_CD'
		'WEIGHT'
	]
	
	for fn in file_names
		lines = fs.readFileSync("../data/#{fn}.txt", {'encoding':'utf8'}).split '\n'
		switch fn
			when 'FOOD_DES'
				try
					db.serialize () ->
						db.run '
							CREATE TABLE FOOD_DES (
								NDB_No 		VARCHAR(5) NOT NULL UNIQUE PRIMARY KEY,
								FdGrp_Cd	VARCHAR(4) NOT NULL,
								Long_Desc	VARCHAR(200) NOT NULL,
								Shrt_Desc	VARCHAR(60) NOT NULL,
								ComName		VARCHAR(100),
								ManufacName	VARCHAR(65),
								Survey		VARCHAR(1),
								Ref_desc	VARCHAR(135),
								Refuse		INTEGER,
								SciName		VARCHAR(65),
								N_Factor	REAL,
								Pro_Factor	REAL,
								Fat_Factor	REAL,
								CHO_Factor	REAL
							)
						'
						stmt = db.prepare 'INSERT INTO FOOD_DES VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)'
						for l in lines
							stmt.run (cleanField(f) for f in l.split('^'))
				catch e
					console.log e
					throw "ERROR READING #{fn}"
					
main()