package data

import (
	"fmt"

	core "github.com/alchemicode/jserv-core"
)

func ProcessQuery(cols []*core.Collection, query *core.Query) []*core.Document {
	var ret []*core.Document
	intersect_ret := func(chunk []*core.Document) {
		if ret != nil {
			ret = Intersect(ret, chunk)
		} else {
			ret = chunk
		}
	}
	if len(query.Collections) > 0 {
		for _, c := range query.Collections {
			if col := FindCollection(cols, c); col != nil {
				intersect_ret(PropogateQuery(ret, col, query))
			}
		}
	} else {
		for _, col := range cols {
			intersect_ret(PropogateQuery(ret, col, query))
		}
	}
	return ret
}

func PropogateQuery(ret []*core.Document, col *core.Collection, query *core.Query) []*core.Document {
	intersect_ret := func(chunk []*core.Document) {
		if ret != nil {
			ret = Intersect[*core.Document](ret, chunk)
		} else {
			ret = Intersect[*core.Document](col.List, chunk)
		}
	}
	for _, v := range query.Has {
		equals_objs := FindHas(col, v)
		intersect_ret(equals_objs)
	}
	for k, v := range query.Equals {
		equals_objs := FindEquals(col, k, v)
		intersect_ret(equals_objs)
	}
	for k, v := range query.NotEquals {
		notequals_objs := FindNotEquals(col, k, v)
		intersect_ret(notequals_objs)
	}
	for k, v := range query.LessThan {
		lt_objs := FindLessThan(col, k, v)
		intersect_ret(lt_objs)
	}
	for k, v := range query.LessThanEqualTo {
		lte_objs := FindLessThanEqual(col, k, v)
		intersect_ret(lte_objs)
	}
	for k, v := range query.GreaterThanEqualTo {
		gte_objs := FindGreaterThanEqual(col, k, v)
		intersect_ret(gte_objs)
	}
	for k, v := range query.GreaterThan {
		gt_objs := FindGreaterThan(col, k, v)
		intersect_ret(gt_objs)
	}
	for k, v := range query.Between {
		vals := v.([]interface{})
		bet_objs := FindBetween(col, k, vals[0], vals[1])
		intersect_ret(bet_objs)
	}
	return ret
}

func ProcessModify(cols []*core.Collection, mod *core.Mod, resp *core.Response) {
	mods := make([]string, 0)
	skips := make([]string, 0)
	if col := FindCollection(cols, mod.Collection); col != nil {
		if doc := FindDocument(col, mod.Document); doc != nil {
			for name, val := range mod.Values {
				if name != "_id" {
					if _, ok := doc.Data[name]; ok {
						doc.Data[name] = val
						mods = append(mods, name)
					} else {
						skips = append(skips, name+": Does not exist in doc_"+mod.Document)
					}
				} else {
					doc.Id = val.(string)
					mods = append(mods, "_d")
				}
			}
			if len(skips) > 0 {
				resp.WithData("ok", "Completed Mod with skips", map[string]interface{}{"mods": mods, "skips": skips})
			} else {
				resp.WithData("ok", "Completed Mod successfully", map[string]interface{}{"mods": mods})
			}
		} else {
			resp.WithoutData("error", fmt.Sprintf("doc_%s does not exist", mod.Document))
		}
	} else {
		resp.WithoutData("error", fmt.Sprintf("Collection %s does not exist", mod.Collection))
	}
}

func ProcessAdd(cols []*core.Collection, add *core.Add, resp *core.Response) {
	adds := make([]string, 0)
	skips := make([]string, 0)
	if col := FindCollection(cols, add.Collection); col != nil {
		for _, doc := range add.Documents {
			if FindDocument(col, doc.Id) == nil {
				adds = append(adds, "doc_"+doc.Id)
				new_doc := new(core.Document)
				new_doc.Id = doc.Id
				new_doc.Data = doc.Data
				col.List = append(col.List, new_doc)
			} else {
				skips = append(skips, fmt.Sprintf("doc_%s: Already exists", doc.Id))
			}
		}
		for k, v := range add.Values {
			if doc := FindDocument(col, k); doc != nil {
				values := v.(map[string]interface{})
				for name, val := range values {
					if name != "_id" {
						if _, ok := doc.Data[name]; !ok {
							doc.Data[name] = val
							adds = append(adds, "doc_"+k+"->"+name)
						} else {
							skips = append(skips, fmt.Sprintf("%s: Already exists in doc_%s", name, doc))
						}
					} else {
						skips = append(skips, k+": \"_id\" is a reserved keyword and cannot be added as a data value\n")
					}
				}
			} else {
				skips = append(skips, fmt.Sprintf("doc_%s: Does not exist", k))
			}
		}
		if len(skips) > 0 {
			resp.WithData("ok", "Completed Add with skips", map[string]interface{}{"adds": adds, "skips": skips})
		} else {
			resp.WithData("ok", "Completed Add successfully", map[string]interface{}{"adds": adds})
		}
	} else {
		resp.WithoutData("error", fmt.Sprintf("Collection %s does not exist", add.Collection))
	}
}

func ProcessDelete(cols []*core.Collection, del *core.Delete, resp *core.Response) {
	mods := make([]string, 0)
	skips := make([]string, 0)
	if col := FindCollection(cols, del.Collection); col != nil {
		if doc := FindDocument(col, del.Document); doc != nil {
			for _, name := range del.Values {
				if name != "_id" {
					if _, ok := doc.Data[name]; ok {
						delete(doc.Data, name)
						mods = append(mods, name)
					} else {
						skips = append(skips, name+": Doesn't exist in doc_"+del.Document)
					}
				} else {
					skips = append(skips, "_id: Unique identifier, cannot be deleted")
				}
			}
			if len(skips) > 0 {
				resp.WithData("ok", "Completed Delete with skips", map[string]interface{}{"deleted": mods, "skips": skips})
			} else {
				resp.WithData("ok", "Completed Delete successfully", map[string]interface{}{"deleted": mods})
			}
		} else {
			resp.WithoutData("error", fmt.Sprintf("doc_%s does not exist", del.Document))
		}
	} else {
		resp.WithoutData("error", fmt.Sprintf("Collection %s does not exist", del.Collection))
	}
}

func ProcessPurge(cols []*core.Collection, del *core.Query, resp *core.Response) {
	var deletes []string
	if len(del.Collections) > 0 {
		for _, c := range del.Collections {
			if col := FindCollection(cols, c); col != nil {
				deletes = PropogatePurge(col, del)
			} else {
				resp.WithoutData("error", fmt.Sprintf("Collection %s does not exist", c))
				return
			}
		}
	} else {
		for _, col := range cols {
			deletes = PropogatePurge(col, del)
		}
	}
	if len(deletes) > 0 {
		resp.WithData("ok", "Successfully purged documents", map[string]interface{}{"purged": deletes})
	} else {
		resp.WithoutData("ok", "No documents deleted")
	}

}

func PropogatePurge(col *core.Collection, del *core.Query) []string {
	to_delete := make([]*core.Document, 0)
	intersect_delete := func(chunk []*core.Document) {
		if len(to_delete) > 0 {
			to_delete = Intersect[*core.Document](to_delete, chunk)
		} else {
			to_delete = Intersect[*core.Document](col.List, chunk)
		}
	}
	for _, v := range del.Has {
		equals_objs := FindHas(col, v)
		intersect_delete(equals_objs)
	}
	for k, v := range del.Equals {
		equals_objs := FindEquals(col, k, v)
		intersect_delete(equals_objs)
	}
	for k, v := range del.NotEquals {
		notequals_objs := FindNotEquals(col, k, v)
		intersect_delete(notequals_objs)
	}
	for k, v := range del.LessThan {
		lt_objs := FindLessThan(col, k, v)
		intersect_delete(lt_objs)
	}
	for k, v := range del.LessThanEqualTo {
		lte_objs := FindLessThanEqual(col, k, v)
		intersect_delete(lte_objs)
	}
	for k, v := range del.GreaterThanEqualTo {
		gte_objs := FindGreaterThanEqual(col, k, v)
		intersect_delete(gte_objs)
	}
	for k, v := range del.GreaterThan {
		gt_objs := FindGreaterThan(col, k, v)
		intersect_delete(gt_objs)
	}
	for k, v := range del.Between {
		vals := v.([]interface{})
		bet_objs := FindBetween(col, k, vals[0], vals[1])
		intersect_delete(bet_objs)
	}
	deletes := make([]string, 0)
	for _, doc := range to_delete {
		deletes = append(deletes, "doc_"+doc.Id)
		col.List = Remove(doc, col.List)
	}
	return deletes
}
