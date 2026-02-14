package datetime

//go:generate goyacc -l -o parse_yacc.go -v parse_yacc.states.txt parse_yacc.y
//go:generate go run ../cmd/glr-generate parse_yacc.y parse_yacc.states.txt
