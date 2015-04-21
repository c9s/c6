package c6

/*
There are 3 scope that users may use interpolation syntax:

 {selector interpolation}  {
	 {property name inpterpolation}: {property value interpolation}
 }

*/
func lexInterpolation(l *Lexer, emit bool) stateFn {
	l.remember()
	var r rune = l.next()
	if r == '#' {
		r = l.next()
		if r == '{' {
			if emit {
				l.emit(T_INTERPOLATION_START)
			}

			r = l.next()
			for r == ' ' {
				r = l.next()
			}
			l.backup()

			r = l.next()
			for r != '}' {
				r = l.next()
			}
			l.backup()

			if emit {
				l.emit(T_INTERPOLATION_INNER)
			}

			l.next() // for '}'
			if emit {
				l.emit(T_INTERPOLATION_END)
			}
			return nil
		}
	}
	l.rollback()
	return nil
}
