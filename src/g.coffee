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

#
#test
#
G = 
	a: ['b']
	b: ['c','d']
	c: ['e']
	d: ['e']
	e: []
console.log topologicalSort G