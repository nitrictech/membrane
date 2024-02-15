
package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/golang/mock/mockgen/model"

	pkg_ "github.com/nitrictech/nitric/cloud/azure/runtime/queue/iface"
)

var output = flag.String("output", "", "The output file name, or empty to use stdout.")

func main() {
	flag.Parse()

	its := []struct{
		sym string
		typ reflect.Type
	}{
		
		{ "AzqueueServiceUrlIface", reflect.TypeOf((*pkg_.AzqueueServiceUrlIface)(nil)).Elem()},
		
		{ "AzqueueQueueUrlIface", reflect.TypeOf((*pkg_.AzqueueQueueUrlIface)(nil)).Elem()},
		
		{ "AzqueueMessageUrlIface", reflect.TypeOf((*pkg_.AzqueueMessageUrlIface)(nil)).Elem()},
		
		{ "DequeueMessagesResponseIface", reflect.TypeOf((*pkg_.DequeueMessagesResponseIface)(nil)).Elem()},
		
		{ "AzqueueMessageIdUrlIface", reflect.TypeOf((*pkg_.AzqueueMessageIdUrlIface)(nil)).Elem()},
		
	}
	pkg := &model.Package{
		// NOTE: This behaves contrary to documented behaviour if the
		// package name is not the final component of the import path.
		// The reflect package doesn't expose the package name, though.
		Name: path.Base("github.com/nitrictech/nitric/cloud/azure/runtime/queue/iface"),
	}

	for _, it := range its {
		intf, err := model.InterfaceFromInterfaceType(it.typ)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Reflection: %v\n", err)
			os.Exit(1)
		}
		intf.Name = it.sym
		pkg.Interfaces = append(pkg.Interfaces, intf)
	}

	outfile := os.Stdout
	if len(*output) != 0 {
		var err error
		outfile, err = os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open output file %q", *output)
		}
		defer func() {
			if err := outfile.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to close output file %q", *output)
				os.Exit(1)
			}
		}()
	}

	if err := gob.NewEncoder(outfile).Encode(pkg); err != nil {
		fmt.Fprintf(os.Stderr, "gob encode: %v\n", err)
		os.Exit(1)
	}
}
