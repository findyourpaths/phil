package glr

//go:generate goyacc -l -o simple_yacc.go -v simple_yacc.states.txt simple_yacc.y
//go:generate go run ../tools/glr/glr_generate.go --in-glr-pkg simple_yacc.y simple_yacc.states.txt
