package main

import (
	"os"
	"github.com/newestuser/contract-committee/app/datastore/mongo"
	log "github.com/cihub/seelog"
	"github.com/newestuser/contract-committee/cmd/config"
	"time"
	"fmt"
	"reflect"
	"unsafe"
	"text/template/parse"
	"github.com/newestuser/contract-committee/app/pars"
)

type Account struct {
	FirstName string
	LastName  string
}

type Purchase struct {
	Date          time.Time
	Description   string
	AmountInCents int
}

type Statement struct {
	FromDate  time.Time
	ToDate    time.Time
	Account   Account
	Purchases []Purchase
}

func main() {

	funcMap := pars.FuncMap{
		"rndStr": rndStr,
		"savе": save,
	}

	template := pars.NewTemplate("email-test").Funcs(funcMap)
	template.Parse("{{$custNum := rndStr}}")
	template.Execute()

	//if err != nil {
	//	panic(err)
	//}


	//funcs := map[string]interface{}{
	//	"rndStr":rndStr,
	//	"save": assertSave("123"),
	//}
	//
	//tree, _ := parse.Parse("foo", `"customerNumber":"{{$custNum := rndStr | save}}"`, "{{", "}}", funcs)
	//
	//for k, v := range tree {
	//	fmt.Println("key:", k)
	//	fmt.Println("val:", v)
	//	fmt.Println("============================================")
	//	root := v.Root
	//	fmt.Println("root:", root)
	//	fmt.Println("============================================")
	//
	//	for i, n := range root.Nodes {
	//		fmt.Println("n", i, ":", n)
	//		fmt.Println("type", n.Type())
	//		printNodeStuff(n)
	//	}
	//}

	//fmt.Println(tree)
	//fmt.Println(err)

}

func printNodeStuff(node parse.Node) {
	switch n := node.(type) {
	case *parse.ActionNode:
		fmt.Println("line", n.Line)
		fmt.Println("pos", n.Pos)
		fmt.Println("pipe", n.Pipe)

		for i, cmd := range n.Pipe.Cmds{
			fmt.Println("pipe cmd",i,cmd)
			fmt.Println("pipe cmd args",i,cmd.Args)
		}

		for i, cmd := range n.Pipe.Decl{
			fmt.Println("pipe decl indent",i,cmd.Ident)
		}
		return
	}
}

func formatAsDollars(valueInCents int) (string, error) {
	dollars := valueInCents / 100
	cents := valueInCents % 100
	return fmt.Sprintf("$%d.%2d", dollars, cents), nil
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%d/%d", day, month, year)
}

func urgentNote(acc Account) string {
	return fmt.Sprintf("You have earned 100 VIP points that can be used for purchases")
}

func save(name, val string) string {
	fmt.Println(val)
	return val
}

func createMockStatement() Statement {
	return Statement{
		FromDate: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
		ToDate: time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
		Account: Account{
			FirstName: "John",
			LastName: "Dow",
		},
		Purchases: []Purchase{
			Purchase{
				Date: time.Date(2016, 1, 3, 0, 0, 0, 0, time.UTC),
				Description: "Shovel",
				AmountInCents: 2326,
			},
			Purchase{
				Date: time.Date(2016, 1, 8, 0, 0, 0, 0, time.UTC),
				Description: "Staple remover",
				AmountInCents: 5432,
			},
		},
	}
}

func initMain() {
	_, err := config.New(os.Args[1:])

	if err != nil {
		log.Critical(err)
		return
	}

	c := mongo.Config{Hosts: []string{"localhost"}, DBName: "committee", Indexes: []*mongo.Index{}}

	database := mongo.NewDatabase(c)

	_ = database
}

func templ() string {
	return ""
}

func change_foo(f *parse.Tree) {
	// Note, simply doing reflect.ValueOf(*f) won't work, need to do this
	pointerVal := reflect.ValueOf(f)
	val := reflect.Indirect(pointerVal)

	member := val.FieldByName("vars")
	ptrToY := unsafe.Pointer(member.UnsafeAddr())
	realPtrToY := (**[]string)(ptrToY)
	*realPtrToY = nil // or &Foo{} or whatever
}

func rndStr() string {
	return "123"
}

func assertSave(wantVal string) (func(string) string) {
	return func(gotVal string) string {

		return gotVal
	}
}
