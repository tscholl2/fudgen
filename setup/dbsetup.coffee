fs = require 'fs'
sqlite3 = require 'sqlite3'
Lazy = require 'lazy'

db = new sqlite3.Database '../data/db'

cleanField = (raw) ->
	if raw.length > 1 and raw[0] == '~' and raw[raw.length-1] == '~'
		return raw.substr 1, raw.length-2
	else
		return raw
		
getPath = (f) ->
	return "../data/#{f}.txt"
		
files = [
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
		
columns =
	DATA_SRC:
		[
			['DataSrc_ID',		'VARCHAR(6) NOT NULL']
			['Authors',			'VARCHAR(255)']
			['Title',			'VARCHAR(255) NOT NULL']
			['Year',			'VARCHAR(4)']
			['Journal',			'VARCHAR(135)']
			['Vol_City',		'VARCHAR(16)']
			['Issue_State',		'VARCHAR(5)']
			['Start_Page',		'VARCHAR(5)']
			['End_Page',		'VARCHAR(5)']
		]
	DATSRCLN:
		[
			['NDB_No',			'VARCHAR(5) NOT NULL']
			['Nutr_No',			'VARCHAR(3) NOT NULL']
			['DataSrc_ID',		'VARCHAR(6) NOT NULL']
		]
	DERIV_CD:
		[
			['Deriv_Cd',		'VARCHAR(4) NOT NULL']
			['Deriv_Desc',		'VARCHAR(120) NOT NULL']
		]
	FD_GROUP:
		[
			['FdGrp_Cd',		'VARCHAR(4) NOT NULL']
			['FdGrp_Desc',		'VARCHAR(60) NOT NULL']
		]
	FOOD_DES:
		[
			['NDB_No', 			'VARCHAR(5) NOT NULL']
			['FdGrp_Cd',		'VARCHAR(4) NOT NULL']
			['Long_Desc',		'VARCHAR(200) NOT NULL']
			['Shrt_Desc',		'VARCHAR(60) NOT NULL']
			['ComName',			'VARCHAR(100)']
			['ManufacName',		'VARCHAR(65)']
			['Survey',			'VARCHAR(1)']
			['Ref_desc',		'VARCHAR(135)']
			['Refuse',			'INTEGER(2)']
			['SciName',			'VARCHAR(65)']
			['N_Factor',		'REAL(4.2)']
			['Pro_Factor',		'REAL(4.2)']
			['Fat_Factor',		'REAL(4.2)']
			['CHO_Factor',		'REAL(4.2)']
		]
	FOOTNOTE:
		[
			['NDB_No',			'VARCHAR(5) NOT NULL']
			['Footnt_No',		'VARCHAR(4) NOT NULL']
			['Footnt_Typ',		'VARCHAR(1) NOT NULL']
			['Nutr_No',			'VARCHAR(3)']
			['Footnt_Txt',		'VARCHAR(200) NOT NULL']
		]
	LANGDESC:
		[
			['Factor_Code',		'VARCHAR(5) NOT NULL']
			['Description',		'VARCHAR(140) NOT NULL']
		]
	LANGUAL:
		[
			['NDB_No',			'VARCHAR(5) NOT NULL']
			['Factor_Code',		'VARCHAR(5) NOT NULL']
		]
	NUT_DATA:
		[
			['NDB_No',			'VARCHAR(5) NOT NULL']
			['Nutr_No',			'VARCHAR(3) NOT NULL']
			['Nutr_Val',		'REAL(10.3) NOT NULL']
			['Num_Data_Pts',	'REAL(5.0) NOT NULL']
			['Std_Error',		'REAL(8.3)']
			['Src_Cd',			'VARCHAR(2) NOT NULL']
			['Deriv_Cd',		'VARCHAR(4)']
			['Ref_NDB_No',		'VARCHAR(5)']
			['Add_Nutr_Mark',	'VARCHAR(1)']
			['Num_Studies',		'INTEGER(2)']
			['Min',				'REAL(10.3)']
			['Max',				'REAL(10.3)']
			['DF',				'INTEGER(4)']
			['Low_EB',			'REAL(10.3)']
			['Up_EB',			'REAL(10.3)']
			['Stat_cmt',		'VARCHAR(10)']
			['AddMod_Date',		'VARCHAR(10)']
			['CC',				'VARCHAR(1)']
		]
	NUTR_DEF:
		[
			['Nutr_No',			'VARCHAR(3) NOT NULL']
			['Units',			'VARCHAR(7) NOT NULL']
			['Tagname',			'VARCHAR(20)']
			['NutrDesc',		'VARCHAR(60) NOT NULL']
			['Num_Dec',			'VARCHAR(1) NOT NULL']
			['SR_Order',		'INTEGER(6) NOT NULL']
		]
	SRC_CD:
		[
			['Src_Cd',			'VARCHAR(2) NOT NULL']
			['SrcCd_Desc',		'VARCHAR(60) NOT NULL']
		]
	WEIGHT:
		[
			['NDB_No',			'VARCHAR(5) NOT NULL']
			['Seq',				'VARCHAR(2) NOT NULL']
			['Amount',			'REAL(5.3) NOT NULL']
			['Msre_Desc',		'VARCHAR(84) NOT NULL']
			['GM_Wgt',			'REAL(7.1) NOT NULL']
			['Num_Data_Pts',	'INTEGER(3)']
			['Std_Dev',			'REAL(7.3)']
		]

primary = 
	WEIGHT:
		['NDB_No','Seq']
	SRC_CD:
		['Src_Cd']
	NUTR_DEF:
		['Nutr_No']
	NUT_DATA:
		['NDB_No','Nutr_No']
	LANGUAL:
		['NDB_No','Factor_Code']
	LANGDESC:
		['Factor_Code']
	FOOTNOTE:
		[]
	FOOD_DES:
		['NDB_No']
	FD_GROUP:
		['FdGrp_Cd']
	DERIV_CD:
		['Deriv_Cd']
	DATSRCLN:
		['NDB_No','Nutr_No','DataSrc_ID']
	DATA_SRC:
		['DataSrc_ID']
	
		
relations = 
	FOOD_DES:
		[
			
		]
	FOOTNOTE:	
		[
			['NDB_No',			'FOOD_DES']
		]
	NUT_DATA:
		[
			['NDB_No',			'FOOD_DES']
		]
	WEIGHT:
		[
			['NDB_No',	 		'FOOD_DES']
		]
	LANGUAL:
		[
			['NDB_No',			'FOOD_DES']
		]
	FD_GROUP:
		[
			['FdGrp_Cd',		'FOOD_DES']
		]
	LANGDESC:
		[
			['Factor_Code',		'LANGUAL']
		]
	NUTR_DEF:
		[
			['Nutr_No',			'NUT_DATA']
		]
	SRC_CD:
		[
			['Src_Cd',			'NUT_DATA']
		]
	DERIV_CD:
		[
			['Deriv_Cd',		'NUT_DATA']
		]
	DATSRCLN:
		[
			['NDB_No',			'NUT_DATA']
			['Nutr_No',			'NUTR_DEF']
		]
	DATA_SRC:
		[
			['DataSrc_ID',		'DATSRCLN']
		]
#MAIN

main = () ->
	#run all db queries in serial, this is sticky so will last
	db.serialize()

	#load tables
	console.log '...loading tables...'
	for f in files
		pk = if primary[f]? then primary[f] else []
		fk = if relations[f]? then relations[f] else []
		sql = "CREATE TABLE #{f}\n\t("
		sql += "\n\t\t" + (c[0]+'\t\t'+c[1] for c in columns[f]).join(',\n\t\t')
		if pk.length > 0
			sql += ",\n\t\tPRIMARY KEY\t(" + pk.join(', ') + ")"
		if fk.length > 0
			sql += ','
		for k in fk
			sql += "\n\t\tFOREIGN KEY(#{k[0]})\tREFERENCES #{k[1]}(#{k[0]})"
		sql += "\n\t)"
		db.run sql
		
	#import files synchronously
	console.log '...importing files...'
	n = files.length
	loadFile = () ->
		if files.length == 0
			return
		#gather file
		f = files.pop()
		console.log "loading file #{f} size=#{fs.statSync(getPath(f)).size}"
		#initialize sql cmds
		stmt = db.prepare "INSERT INTO #{f} VALUES (#{('?' for i in columns[f]).join(',')})"
		lazy = new Lazy fs.createReadStream getPath(f), {encoding:'utf8'}
		lazy
		.on 'pipe', () ->
			return loadFile()
		.lines
		.forEach (l) ->
			if l.length > 0
				#insert data
				stmt.run (cleanField(f) for f in l.toString().split('^'))
	#start loads
	loadFile()

main()