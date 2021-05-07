package gosqlgo

import "testing"

type User struct {
	Name     string `db:"name"`
	Lastname string `db:"lastname"`
}

func TestNormalSelect(t *testing.T) {
	want := "SELECT name, lastname FROM users"
	users := AddTable("users")
	var persona User
	users.Columns(persona)
	if got := users.GetQuery(); got != want {
		t.Errorf("GetQuery() = %q, want %q", got, want)
	}
}

func TestWhereSelect(t *testing.T) {
	want := "SELECT name, lastname FROM users WHERE idUser = '1' AND name = 'Andres'"
	users := AddTable("users")
	var persona User
	users.Columns(persona)
	users.Where(map[string]string{"idUser": "1", "name": "Andres"})
	if got := users.GetQuery(); got != want {
		t.Errorf("GetQuery() = %q, want %q", got, want)
	}
}

func TestOrder(t *testing.T) {
	want := "SELECT * FROM users ORDER BY firstname, lastname"
	users := AddTable("users")
	users.Order([]string{"firstname", "lastname"})
	if got := users.GetQuery(); got != want {
		t.Errorf("GetQuery() = %q, want %q", got, want)
	}
}

func TestGroup(t *testing.T) {
	want := "SELECT * FROM users GROUP BY firstname"
	users := AddTable("users")
	users.Group([]string{"firstname"})
	if got := users.GetQuery(); got != want {
		t.Errorf("GetQuery() = %q, want %q", got, want)
	}
}

func TestFetchAll(t *testing.T) {
	want := "SELECT * FROM users"
	users := AddTable("users")
	if got := users.Select(); got != want {
		t.Errorf("GetQuery() = %q, want %q", got, want)
	}
}

func TestJoin(t *testing.T) {
	want := "SELECT habitaciones, baños FROM users INNER JOIN casa ON users.idCasa = casa.idCasa"
	users := AddTable("users")
	users.Join("casa", "users.idCasa = casa.idCasa", []string{"habitaciones", "baños"}, users.InnerJoin())
	if got := users.GetQuery(); got != want {
		t.Errorf("GetQuery() = %q, want %q", got, want)
	}
}

func TestMultiJoin(t *testing.T) {
	want := "SELECT * FROM users LEFT JOIN casa ON users.idCasa = casa.idCasa RIGHT JOIN oficina ON users.idOficina = oficina.idOficina"
	users := AddTable("users")
	users.Join("casa", "users.idCasa = casa.idCasa", nil, users.LeftJoin())
	users.Join("oficina", "users.idOficina = oficina.idOficina", nil, users.RightJoin())
	if got := users.GetQuery(); got != want {
		t.Errorf("GetQuery() = %q, want %q", got, want)
	}
}

func TestInsert(t *testing.T) {
	want := "INSERT INTO users (firstname, lastname) VALUES('Andres', 'Castro');"
	user := AddTable("users")
	if got := user.Insert(map[string]string{"firstname": "Andres", "lastname": "Castro"}); got != want {
		t.Errorf("Insert() = %q, want %q", got, want)
	}
}

func TestDelete(t *testing.T) {
	want := "DELETE FROM users WHERE idUser = '1' AND firstname = 'Andres'"
	user := AddTable("users")
	if got := user.Delete(map[string]string{"idUser": "1", "firstname": "Andres"}); got != want {
		t.Errorf("Delete() = %q, want %q", got, want)
	}
}
func TestUpdate(t *testing.T) {
	want := "UPDATE users SET firstname = 'Andres', lastname = 'Castro' WHERE idUser = '1'"
	user := AddTable("users")
	if got := user.Update(map[string]string{"firstname": "Andres", "lastname": "Castro"}, map[string]string{"idUser": "1"}); got != want {
		t.Errorf("Update() = %q, want %q", got, want)
	}
}
