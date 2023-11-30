package data

import (
	core "github.com/alchemicode/jserv-core"
)

// Checks for a collection of the given name
func FindCollection(c []*core.Collection, name string) *core.Collection {
	for _, v := range c {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// Checks for an object of the given id in a collection
func FindDocument(c *core.Collection, id string) *core.Document {
	for _, v := range c.List {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func RemoveDocument(c *core.Collection, id string) {
	for i, v := range c.List {
		if v.Id == id {
			c.List[i] = c.List[len(c.List)-1] // Copy last element to index i.
			c.List[len(c.List)-1] = nil       // Erase last element (write zero value).
			c.List = c.List[:len(c.List)-1]
			break
		}
	}
}

func RemoveValue(c *core.Collection, id string, val string) {
	for _, v := range c.List {
		if v.Id == id {
			for j := range v.Data {
				if j == val {
					delete(v.Data, j)
					break
				}
			}
		}
	}
}

// Checks for objects of a given attribute in a collection
func FindHas(c *core.Collection, key string) []*core.Document {
	data := make([]*core.Document, 0)
	for _, d := range c.List {
		if _, ok := d.Data[key]; ok {
			data = append(data, d)
		}
	}
	return data
}

// Checks for objects of a given attribute in a collection
func FindEquals(c *core.Collection, att string, val interface{}) []*core.Document {
	data := make([]*core.Document, 0)
	for _, d := range c.List {
		if att == "_id" {
			if d.Id == val {
				data = append(data, d)
			}
		} else {
			if v, ok := d.Data[att]; ok {
				if v == val {
					data = append(data, d)
				}
			}
		}
	}
	return data
}

// Checks for objects of a given attribute in a collection
func FindNotEquals(c *core.Collection, att string, val interface{}) []*core.Document {
	data := make([]*core.Document, 0)
	for _, d := range c.List {
		if att == "_id" {
			if d.Id != val {
				data = append(data, d)
			}
		} else {
			if v, ok := d.Data[att]; ok {
				if v != val {
					data = append(data, d)
				}
			}
		}
	}
	return data
}

// Checks for objects of a given attribute in a collection
func FindLessThan(c *core.Collection, att string, val interface{}) []*core.Document {
	data := make([]*core.Document, 0)
	for _, d := range c.List {
		if att == "_id" {
			if d.Id < val.(string) {
				data = append(data, d)
			}
		} else {
			if v, ok := d.Data[att]; ok {
				if comp, ok := TypeSwitchComparison(v, val); ok {
					if comp == -1 {
						data = append(data, d)
					}
				}

			}
		}
	}
	return data
}

// Checks for objects of a given attribute in a collection
func FindLessThanEqual(c *core.Collection, att string, val interface{}) []*core.Document {
	data := make([]*core.Document, 0)
	for _, d := range c.List {
		if att == "_id" {
			if d.Id < val.(string) {
				data = append(data, d)
			}
		} else {
			if v, ok := d.Data[att]; ok {
				if comp, ok := TypeSwitchComparison(v, val); ok {
					if comp == -1 || comp == 0 {
						data = append(data, d)
					}
				}

			}
		}
	}
	return data
}

// Checks for objects of a given attribute in a collection
func FindGreaterThanEqual(c *core.Collection, att string, val interface{}) []*core.Document {
	data := make([]*core.Document, 0)
	for _, d := range c.List {
		if att == "_id" {
			if d.Id < val.(string) {
				data = append(data, d)
			}
		} else {
			if v, ok := d.Data[att]; ok {
				if comp, ok := TypeSwitchComparison(v, val); ok {
					if comp == 1 || comp == 0 {
						data = append(data, d)
					}
				}

			}
		}
	}
	return data
}

// Checks for objects of a given attribute in a collection
func FindGreaterThan(c *core.Collection, att string, val interface{}) []*core.Document {
	data := make([]*core.Document, 0)
	for _, d := range c.List {
		if att == "_id" {
			if d.Id < val.(string) {
				data = append(data, d)
			}
		} else {
			if v, ok := d.Data[att]; ok {
				if comp, ok := TypeSwitchComparison(v, val); ok {
					if comp == 1 {
						data = append(data, d)
					}
				}

			}
		}
	}
	return data
}

// Checks for objects of a given attribute in a collection
func FindBetween(c *core.Collection, att string, val1 interface{}, val2 interface{}) []*core.Document {
	data := make([]*core.Document, 0)
	for _, d := range c.List {
		if att == "_id" {
			if d.Id >= val1.(string) && d.Id <= val2.(string) {
				data = append(data, d)
			}
		} else {
			if v, ok := d.Data[att]; ok {
				if comp, ok := TypeSwitchComparison(v, val1); ok {
					if comp2, ok2 := TypeSwitchComparison(v, val2); ok2 {
						if (comp == 1 || comp == 0) && (comp2 == -1 || comp2 == 0) {
							data = append(data, d)
						}
					}
				}
			}
		}
	}
	return data
}
