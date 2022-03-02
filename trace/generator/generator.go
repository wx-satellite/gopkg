package generator


type Generator interface {
	Generate(file string)([]byte,error)

}