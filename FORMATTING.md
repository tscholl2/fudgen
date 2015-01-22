#Recipe Formatting

My own version of how recipes should be formatted.

###The Idea

The formatting is yaml-inspired and based on expressions
and makefiles. The idea is a recipe is a collection of
operations that require some input, either the output of
previous operations or bare ingredients. For example
`spaghetti` should require `pasta` and `sauce` before
it's done. The operation would be `mix`. So you might say 
```
spaghetti = mix (pasta, sauce)
```
or
```
spaghetti = pasta + sauce
```
This is the idea.

### The Actual Format

The actual formatting would be for the examples above
would be
```
spaghetti
	+mix
		.5 min
		pasta
			.1 box
		sauce
			.15 oz
```
Or maybe you want spaghetti from scratch
```
spaghetti
	+mix
		.time 1min
		pasta
			+boil
				.8 min
				.description 'Boil till aldente'
				spaghetti noodles
					.1 box
		sauce
			+saute
				.25 min
				.setting medium
				+smash
					.10 min
					tomatoes
						.1 lb
				+mix
					.10 min
					oregano
						.4 tbl
					salt
						.1 tbl
					pepper
						.1 tbl
```
From this it's easy to see recipes are built out of
operations and ingrediants with specified attributes.

## #Operations

Operations are denoted on a new line by one of the key
words underneath them should be all the parameters as
well as the input ingredients, which can either be
the output of other operations or bare ingredients.
Key points:
* All operations are prefixed with a `+`.
* All operations must have a .time attribute
(see attributes below)
* Some operations require extra attributes

####List
Here are the operations in the format

`+operation`: `.attribute`,`[.optional_attribute]`,...

Note that the `.time` attribute is required, so it is best
to put it first and then it doesn't need the '.time' prefix
(see examples above).

- `+bake`: `.time` (1m/30s/2h), `.temp` (450f/300c)
- `+broil`: `.time`, `.setting` (low/medium/high)
- `+mix`: `.time`, `[.tool]` (whisk/smash/hands/light),`[.container]` (bucket)
- `+saute`: `.time`, `.setting`
- `+fry`: `.time`, `.setting`
- `+cut`: `.time`, `[.pieces]` (3/5/119)
- `+slice`: `.time`, `[.size]` (2in), `[.pieces]`
- `+chop`: `.time`, `[.size]`
- `+drain`: `.time`
- `+grate`: `.time`, `[.size]`
- `+nuke`: `.time`,`.setting`

## #Ingredients
Ingredients are denoted on a new line by the plain
description of the item. The `.quantity` attribute
follows same rules as the `.time` for operations.

Ingrediants can also have the `.random` attribute
which should be a number between `0` and `100`.
This represents how this ingredient should be
randomized on recipe generation. A `0` should give
fairly consistent results. A `100` means `onion` might
turn into `raw steak`.

```
ingrediant_description
	.quantity
	.attribute
	.attribute
	[.random]
```

## #Attributes
Self explanatory. The main thing is probably the unit
conventions. Some examples are below. See the parser for
any other questions.

`.time`
- `2m`: 2 minutes
- `2s`: 2 seconds
- `2h`: 2 hours
- `2d`: 2 days
- `2w`: 2 weeks
- `2h 2m`: 2 hours and 2 minutes
- `3.7h`: 3.7 hours

`.temp`
- `250`: 250 degrees Fahrenheit
- `250f`: same as `250`
- `250c`: 250 degrees Celsius

`.setting`
- `low`: low
- `medium`: medium
- `high`: high

`.size` and `.quantity`
- `1in^3`: 1 inch cube
- `2cm`: 2 centimetres
- `7mm`: 7 millimetres
- `1tbl`: one tablespoon

