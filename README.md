# Toy Factory

toyorm generate tool 

- use to sql table to create Model's bind struct

### Install

    go get -u github.com/bigpigeon/toy-factory

### Go version

    go-1.10
    go-1.9
    
### Example

goto your go project directory

    cd $GOPATH/you/project/
    toy-factory -name mysql -url "root:@tcp(localhost:3306)/toyorm?charset=utf8&parseTime=True"
    
and then was generate toy_models.go file in your package