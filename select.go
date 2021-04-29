package gosqlgo

import "strings"

type Table struct {
	table  string
	fields []string
	where  map[string]string
	order  []string
	joins  map[string]joinTable
	group  []string
	offset int
	limit  int
}

type joinTable struct {
	table    string
	on       string
	fields   []string
	typeJoin int
}

func AddTable(table string) *Table {
	p := Table{table: table}
	p.fields = nil
	p.group = nil
	p.joins = make(map[string]joinTable)
	return &p
}

func (r *Table) InnerJoin() int {
	return 1
}

func (r *Table) LeftJoin() int {
	return 2
}

func (r *Table) RightJoin() int {
	return 3
}

func (r *Table) OuterJoin() int {
	return 4
}

func (r *Table) Select() string {
	query := "SELECT * FROM " + r.table
	return query
}

func (r *Table) Insert(newdata map[string]string) string {
	fields := make([]string, 0, len(newdata))
	data := make([]string, 0, len(newdata))
	if len(newdata) != 0 {
		for key, value := range newdata {
			fields = append(fields, key)
			data = append(data, value)
		}
	}
	insert := "INSERT INTO " + r.table + " (" + strings.Join(fields, ", ") + ") VALUES('" + strings.Join(data, "', '") + "');"
	return insert
}

func (r *Table) Update(changes map[string]string, conditions map[string]string) string {
	values := make([]string, 0, len(changes))
	if len(changes) != 0 {
		for field, data := range changes {
			values = append(values, field+" = '"+data+"'")
		}
	}
	where := make([]string, 0, len(conditions))
	if len(conditions) != 0 {
		for field, data := range conditions {
			where = append(where, field+" = '"+data+"'")
		}
	}
	update := "UPDATE " + r.table + " SET " + strings.Join(values, ", ") + " WHERE " + strings.Join(where, " AND ")
	return update
}

func (r *Table) Delete(conditions map[string]string) string {
	where := make([]string, 0, len(conditions))
	if len(conditions) != 0 {
		for field, data := range conditions {
			where = append(where, field+" = '"+data+"'")
		}
	}
	delete := "DELETE FROM " + r.table + " WHERE " + strings.Join(where, " AND ")
	return delete
}

func (r *Table) Columns(fields []string) {
	r.fields = fields
}

func (r *Table) Expression(field string) {
	r.fields = append(r.fields, field)
}

func (r *Table) Where(where map[string]string) {
	r.where = where
}

func (r *Table) Order(order []string) {
	r.order = order
}

func (r *Table) Group(group []string) {
	r.group = group
}

func (r *Table) Limit(offset int, limit int) {
	r.offset = offset
	r.limit = limit
}

func (r *Table) GetQuery() string {
	join := ""
	fieldsJoin := ""
	if len(r.joins) != 0 {
		for njoin, data := range r.joins {
			inner := " "
			if data.typeJoin == 2 {
				inner += "LEFT"
			} else if data.typeJoin == 3 {
				inner += "RIGHT"
			} else if data.typeJoin == 4 {
				inner += "OUTER"
			} else {
				inner += "INNER"
			}
			join += inner + " JOIN " + njoin + " ON " + data.on
			coma := ""
			if r.fields != nil {
				coma += ", "
			}
			if data.fields != nil {
				fieldsJoin += coma + strings.Join(data.fields, ", ")
			}
		}
	}
	fields := ""
	if len(r.fields) != 0 {
		fields += strings.Join(r.fields, ", ")
	} else {
		if fieldsJoin == "" {
			fields += "*"
		}
	}
	whereString := ""
	if len(r.where) != 0 {
		where := make([]string, 0, len(r.where))
		for field, data := range r.where {
			where = append(where, field+" = '"+data+"'")
		}
		whereString = " WHERE " + strings.Join(where, " AND ")
	}
	order := ""
	if len(r.order) != 0 {
		order += " ORDER BY " + strings.Join(r.order, ", ")
	}
	group := ""
	if r.group != nil {
		group += " GROUP BY " + strings.Join(r.group, ", ")
	}
	query := "SELECT " + fields + fieldsJoin + " FROM " + r.table + join + whereString + order + group
	return query
}

func (r *Table) Join(table string, on string, fields []string, typeJoin int) {
	join := joinTable{table: table, on: on, fields: fields, typeJoin: typeJoin}
	r.joins[table] = join
}
