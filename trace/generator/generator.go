package generator


type Generator interface {
	Generate(string)([]byte,error)

}