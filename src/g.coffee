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

timePast = (workers) ->
	return min (sum(w) for w in workers)

class Worker
	constructor: ->
		this.history = []
		this.job = null
	finish: ->
		j = this.job
		this.job = null
		this.history.push j
		return j

task = (V,E,n=1) ->
	#initialize graph
	G = {}
	for v of V
		G[v] =
			in: []
			out: []
	for e in E
		G[e[1]].in.push e[0]
		G[e[0]].out.push e[1]
	#initialize workers
	workers = (new Worker() for i in [1..n])
	workerSort = (a,b) ->
		if not a.job?
			return -1
		if not b.job?
			return 1
		return V[a.job] < V[b.job]
	#run scheduling
	while not isEmpty G
		#next worker finishes
		v = workers[0].finish()		
		#remove task
		for w in G[v].out
			G[w].in.splice G[w].in.indexOf(v), 1
		delete G[v]	
		#reorder workers
		workers.sort workerSort
		#find starting nodes
		S = (v for v,n of G when n.in.length == 0)
		#find available workers
		W = (w for w in workers when w.job == null)
		if S.length == 0
			throw 'Cycle?!'
		while S.length > 0 and W.length > 0
			#assign task
			v = S.pop()
			w = W.pop()
			w.job = v
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
