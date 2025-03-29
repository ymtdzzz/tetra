package adapter

import (
	"sort"

	"github.com/ymtdzzz/tetra/config"
)

type DBConnection struct {
	Name     string
	DBConfig config.DBConfig
	Adapter  Adapter
}

type DBConnections []*DBConnection

func (c DBConnections) Close() error {
	for _, conn := range c {
		if err := conn.Adapter.Close(); err != nil {
			return err
		}
	}
	return nil
}

func NewDBConnections(config config.Config) DBConnections {
	type kv struct {
		key string
		t   string
	}

	var sorted []kv
	for k, v := range config {
		sorted = append(sorted, kv{key: k, t: string(v.Type)})
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].t == sorted[j].t {
			return sorted[i].key < sorted[j].key
		}
		return sorted[i].t < sorted[j].t
	})

	var connections []*DBConnection
	for _, s := range sorted {
		switch s.t {
		case "mysql":
			adapter, err := NewMySQLAdapter(&MySQLConfig{
				Host:     config[s.key].Host,
				Port:     config[s.key].Port,
				Username: config[s.key].User,
				Password: config[s.key].Password,
			})
			if err != nil {
				continue
			}

			connections = append(connections, &DBConnection{
				Name:     s.key,
				DBConfig: config[s.key],
				Adapter:  adapter,
			})
		}
	}

	return connections
}
