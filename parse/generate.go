package parse

//go:generate goyacc -l -o datetime_yacc.go -v datetime_yacc.states.txt datetime_yacc.y
/* //go:generate goyacc -o parse_datetime_ranges.go -p "retrieve" parse_datetime_ranges.y */
