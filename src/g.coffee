#
#http://en.wikipedia.org/wiki/Topological_sorting
#
topologicalSort = (G) ->
	
	#build backwards graph
	Grev = {}
	for v of G
		Grev[v] = []
	for v1 of G
		for v2 in G[v1]
			Grev[v2].push v1

	##topological sort##
			
	S = (v for v,n of Grev when n.length == 0) #vertices with no incoming edges
	L = [] #will contain sorted elements
	while S.length > 0
		n = S.pop() #remove node from S
		L.push n #add n to tail of L
		for m in (m for m in G[n]) #for each node m with an edge e from n to m do
			#remove edge e from the graph
			G[n].splice G[n].indexOf(m), 1
			Grev[m].splice Grev[m].indexOf(n), 1
			if Grev[m].length == 0 #if m has no other incoming edges then
				S.push m #insert m into S
	if Math.max((n.length for v,n of G)) > 0 #if graph has edges then
		throw 'Error, graph has cycle' #return error (graph has at least one cycle)
	else
		return L #return L (a topologically sorted order)

isEmpty = (o) ->
	if not o?
		return true
	for k,v of o
		return false
	return true

min = (arr) ->
	m = arr[0]
	for a in arr
		if a < m
			m = a
	return m

task = (V,E,n=1) ->
	E = (e for e in E) #copy input
	Grev = {}
	for v of V
		Grev[v] = []
	for e in E
		Grev[e[1]].push e[0]
	workers = ([] for i in [1..n])
	totals = (0 for i in [1..n])
	while not isEmpty Grev
		#find doable nodes
		S = (v for v,n of Grev when n.length == 0)
		if S.length == 0
			throw 'Cycle?!'
		while S.length > 0
			v = S.pop()
			i = totals.indexOf min totals
			workers[i].push v
			totals[i] += V[v]
			delete Grev[v]
			d = (e for e in E when e[0] == v)
			for e in d
				Grev[e[1]].splice Grev[e[1]].indexOf(e[0]), 1
				E.splice E.indexOf(e), 1
	return workers
#
#test
#
V = 
	'a':10
	'b':10
	'c':10
	'd':10
	'e':10
E = [
	['a','b']
	['b','c']
	['b','d']
	['c','e']
	['d','e']
]

console.log task V,E,2

G = 
	a: ['b']
	b: ['c','d']
	c: ['e']
	d: ['e']
	e: []
#console.log topologicalSort G
