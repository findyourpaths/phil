package datetime

//go:generate goyacc -l -o datetime_yacc.go -v datetime_yacc.states.txt datetime_yacc.y
//go:generate go run ../tools/glr/glr_generate.go datetime_yacc.y datetime_yacc.states.txt
