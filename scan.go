package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type scanner struct {
	Root    string
	Results []*moduleReference
	Paths   []string
}

const terraformSourceFileExt = ".tf"

func (s *scanner) ScanDir(path string) error {
	module, err := tfconfig.LoadModule(path)
	if err != nil {
		return fmt.Errorf("read terraform module %q: %v", path, err)
	}
	var modules []*moduleReference
	for _, m := range module.ModuleCalls {
		modules = append(modules, &moduleReference{
			Name:    m.Name,
			Source:  m.Source,
			Version: &m.Version,
			Path:    m.Pos.Filename,
		})
	}
	for _, m := range modules {
		if err := m.ParseSource(); err != nil {
			log.Printf("parse module source: %v", err)
		}
		s.Results = append(s.Results, m)
		s.Paths = append(s.Paths, path)
	}
	return nil
}
