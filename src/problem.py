import copy
import random
import sympy

def weighted_choice(choices):
	total = sum(w for c, w in choices)
	r = random.uniform(0, total)
	upto = 0
	for c, w in choices:
		if upto + w > r:
			return c
		upto += w
	assert False, "Shouldn't get here"

def problem_creator(symbol=None, _choices=None):
	if symbol is None:
		symbol = sympy.symbols('x')

	if _choices is None:
		_choices = copy.copy(problem_creator.choices)
		outer = weighted_choice(_choices.items())(symbol)
		while not isinstance(outer, Binary):
			outer = weighted_choice(_choices.items())(symbol)
	else:
		outer = weighted_choice(_choices.items())(symbol)

	if isinstance(outer, Binary):
		for _class in _choices:
			if issubclass(_class, Binary):
				_choices[_class] /= 4

	outer.args = [problem_creator(symbol, _choices) for _ in outer.args]
	return outer

class Expression:
	__slots__ = ('args', 'symb')

	def __init__(self, symb):
		self.args = []
		self.symb = symb

class Unary(Expression):
	def __init__(self, *args, **kwds):
		super().__init__(*args, **kwds)
		self.args = [None]

class Binary(Expression):
	def __init__(self, *args, **kwds):
		super().__init__(*args, **kwds)
		self.args = [None, None]

class Add(Binary):
	def to_sympy(self):
		return self.args[0].to_sympy()+self.args[1].to_sympy()

	def tags(self):
		return self.args[0].tags().union(self.args[1].tags())

class Sub(Binary):
	def to_sympy(self):
		return self.args[0].to_sympy()-self.args[1].to_sympy()

	def tags(self):
		return self.args[0].tags().union(self.args[1].tags())

class Mul(Binary):
	def to_sympy(self):
		return self.args[0].to_sympy()*self.args[1].to_sympy()

	def tags(self):
		return frozenset(('product_rule',)).union(
				self.args[0].tags().union(self.args[1].tags()))

class Div(Binary):
	def to_sympy(self):
		return self.args[0].to_sympy()/self.args[1].to_sympy()

	def tags(self):
		return frozenset(('quotient_rule',)).union(
				self.args[0].tags().union(self.args[1].tags()))

class ConstMul(Unary):
	__slots__ = ('const',)

	def __init__(self, *args, **kwds):
		super().__init__(*args, **kwds)
		choices = [random.randint(2, 9), random.randint(-9, -1)]
		self.const = random.choice(choices)

	def to_sympy(self):
		return self.const*self.args[0].to_sympy()

	def tags(self):
		return self.args[0].tags()

class Monomial(Expression):
	__slots__ = ('exp',)

	def __init__(self, *args, **kwds):
		self.exp = random.randint(2, 9)
		super().__init__(*args, **kwds)

	def to_sympy(self):
		return self.symb**self.exp

	def tags(self):
		return frozenset(('power_rule',))

class Polynomial(Expression):
	__slots__ = ('expr',)

	def __init__(self, *args, **kwds):
		degree = kwds.pop('degree', random.randint(3, 9))
		super().__init__(*args, **kwds)
		self.expr = sympy.random_poly(self.symb, degree, -9, 9)

	def to_sympy(self):
		return self.expr

	def tags(self):
		return frozenset(('power_rule',))

class Quadratic(Polynomial):
	def __init__(self, *args, **kwds):
		super().__init__(*args, degree=2, **kwds)

class Linear(Polynomial):
	def __init__(self, *args, **kwds):
		super().__init__(*args, degree=1, **kwds)

	def tags(self):
		return frozenset()

class Root(Expression):
	__slots__ = ('root',)

	def __init__(self, *args, **kwds):
		self.root = random.randint(2, 9)
		super().__init__(*args, **kwds)

	def to_sympy(self):
		return sympy.root(self.symb, self.root)

	def tags(self):
		return frozenset(('power_rule',))

class ExpOffBase(Expression):
	__slots__ = ('base',)

	def __init__(self, *args, **kwds):
		super().__init__(*args, **kwds)
		self.base = random.randint(2, 9)

	def to_sympy(self):
		return self.base**self.symb

	def tags(self):
		return frozenset(('exponential',))

class Exp(Expression):
	def to_sympy(self):
		return sympy.exp(self.symb)

	def tags(self):
		return frozenset(('exponential',))

class ExpMon(Binary):
	def to_sympy(self):
		return self.args[0].to_sympy()**self.args[1].to_sympy()

	def tags(self):
		return frozenset(('exponential',)).union(
				self.args[0].tags().union(self.args[1].tags()))

class Log(Expression):
	def to_sympy(self):
		return sympy.log(self.symb)

	def tags(self):
		return frozenset(('logarithmic',))

class Trig(Expression):
	def tags(self):
		return frozenset(('trigonometric',))

class Sin(Trig):
	def to_sympy(self):
		return sympy.sin(self.symb)

class Cos(Trig):
	def to_sympy(self):
		return sympy.cos(self.symb)

class Tan(Trig):
	def to_sympy(self):
		return sympy.tan(self.symb)

class Sec(Trig):
	def to_sympy(self):
		return sympy.sec(self.symb)

class Csc(Trig):
	def to_sympy(self):
		return sympy.csc(self.symb)

class Cot(Trig):
	def to_sympy(self):
		return sympy.cot(self.symb)

class Comp(Binary):
	def to_sympy(self):
		return self.args[0].to_sympy().subs(self.symb, self.args[1].to_sympy())

	def tags(self):
		return frozenset(('chain_rule',)).union(
				self.args[0].tags().union(self.args[1].tags()))


problem_creator.choices = {
		Add			: 16,
		#Sub			: 1,
		Mul			: 8,
		Div			: 6,
		Comp	   	: 8,
		ConstMul   	: 16,
		Monomial   	: 4,
		Quadratic   : 4,
		Linear	  : 8,
		Root	   	: 2,
		ExpOffBase 	: 1,
		Exp			: 4,
		ExpMon	 	: 1,
		Log			: 2,
		Sin			: 8,
		Cos			: 8,
		Tan			: 4,
		Sec			: 1,
		Csc			: 1,
		Cot			: 1,
		}


def sympy_to_mathjs(expr):
	if isinstance(expr, sympy.Mul):
		return '('+'*'.join(map(sympy_to_mathjs, expr.args))+')'
	elif isinstance(expr, sympy.Add):
		return '('+'+'.join(map(sympy_to_mathjs, expr.args))+')'
	elif isinstance(expr, sympy.Pow):
		return '('+'^'.join(map(sympy_to_mathjs, expr.args))+')'
	elif isinstance(expr, (sympy.sin, sympy.cos, sympy.tan, sympy.sec, sympy.csc, sympy.cot, sympy.exp, sympy.log)):
		return expr.__class__.__name__+'('+','.join(map(sympy_to_mathjs, expr.args))+')'
	elif isinstance(expr, sympy.Symbol):
		return expr.name
	elif isinstance(expr, sympy.Number):
		if expr < 0:
			return '(' + str(expr) + ')'
		else:
			return str(expr)
	assert False, expr

