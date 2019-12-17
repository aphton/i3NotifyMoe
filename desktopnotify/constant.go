package desktopnotify

// JadeTemplate template for notify-send
const JadeTemplate = `
if len(Released) + len(Today) + len(Unreleased) == 0
	b
		span(face='monospace', font='16px', color='#ff5050')
			| User is not watching\n
			| any currently airing\n
			| anime!
else
	if len(Released) > 0
		b
			span(face='monospace', font='20px', color='#50ff50')
				| released
		b
			span(face='monospace', font='16px')
				each $val in Released
					#{"\n    " + $val}
		br
		br
	if len(Today) > 0
		b
			span(face='monospace', font='20px', color='#ffff50')
				| today
		b
			span(face='monospace', font='16px')
				each $val in Today
					#{"\n      " + $val}
		br
		br
	if len(Unreleased) > 0
		b
			span(face='monospace', font='20px', color='#ff5050')
				| unreleased
		b
			span(face='monospace', font='16px')
				each $val in Unreleased
					#{"\n    " + $val}
`
