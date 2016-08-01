package models

import (
	"github.com/nullbio/sqlboiler/boil"
	"github.com/nullbio/sqlboiler/boil/qm"
)

// M type is for providing where filters to Where helpers.
type M map[string]interface{}

// NewQueryG initializes a new Query using the passed in QueryMods
func NewQueryG(mods ...qm.QueryMod) *boil.Query {
	return NewQuery(boil.GetDB(), mods...)
}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(exec boil.Executor, mods ...qm.QueryMod) *boil.Query {
	q := &boil.Query{}
	boil.SetExecutor(q, exec)
	qm.Apply(q, mods...)

	return q
}

