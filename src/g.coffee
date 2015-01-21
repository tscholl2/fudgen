isEmpty = (o) ->
	if not o?
		return true
	for k,v of o
		return false
	return true

class Worker
	constructor: ->
		@job = null
		@time = 0
		@history = []
	assign: (j,time)->
		@job = j
		@time = time
	finish: ->
		j = @job
		t = @time
		@job = null
		@time = 0
		return j
	work: (t)->
		@time = Math.max @time-t,0
		if @history.length > 0 and @history[@history.length-1].job == @job
			@history[@history.length-1].time += t
		else
			@history.push {job:@job,time:t}
	toString: ->
		return {job: @job,time: @time}

class WorkerGroup
	constructor: (n) ->
		@workers = (new Worker() for i in [1..n])
	available: ->
		return (w for w in @workers when not w.job?)
	next: ->
		t = Math.min (w.time for w in @workers when w.job?)...
		for w in @workers
			w.work t
		done = []
		for w in @workers when w.time == 0 and w.job?
			done.push w.finish()
		return done
	current: (j) ->
		for w in @workers
			if w.job == j
				return true
		return false
	busy: ->
		for w in @workers
			if w.job?
				return true
		return false
	schedule: ->
		return (w.history for w in @workers)
	toString: ->
		return (w.toString() for w in @workers)

task = (V,E,n=1) ->
	# initialize graph
	G = {}
	for v of V
		G[v] =
			in: []
			out: []
	for e in E
		G[e[1]].in.push e[0]
		G[e[0]].out.push e[1]
	# initialize workers
	workers = new WorkerGroup n
	# run scheduling
	while not isEmpty G
		# console.log 'starting loop'
		# console.log workers.toString()
		# find starting nodes that aren't being worked on
		S = (v for v,n of G when n.in.length == 0 and not workers.current v)
		# find available workers
		W = workers.available()
		if S.length == 0 and not workers.busy()
			throw 'Cycle?!'
		while S.length > 0 and W.length > 0
			# assign task
			v = S.pop()
			w = W.pop()
			w.assign v, V[v]
		# console.log 'mid loop'
		# console.log workers.toString()
		# finish any jobs
		F = workers.next()
		# remove tasks
		for v in F
			for w in G[v].out
				G[w].in.splice G[w].in.indexOf(v), 1
			delete G[v]	
		# console.log 'finished loop'
		# console.log workers.toString()
	# console.log workers.toString()
	s = workers.schedule()
	s.reverse()
	return s
# 
# test
# 
V = 
	'a':5
	'b':10
	'c':6
	'd':10
	'e':4
E = [
	['a','b']
	['b','c']
	['b','d']
	['c','e']
	['d','e']
]

console.log task V,E,2
