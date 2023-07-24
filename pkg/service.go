package pkg

/*
CREATE TABLE IF NOT EXISTS `employee` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `salary` double NOT NULL,
  `age` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 AUTO_INCREMENT=158;
*/
import (
	"fmt"
)

type Employee struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    string `json:"age"`
}

func SaveEmployee(emp *Employee) error {
	sql := "INSERT INTO employees(name, age, salary) VALUES( ?, ?, ?)"
	stmt, err := MysqlDB.Prepare(sql)

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(emp.Name, emp.Salary, emp.Age)
	if err != nil {
		return err
	}
	fmt.Println(result.LastInsertId())
	return nil
}

func DeleteEmployee(empId string) error {
	sql := "Delete FROM employees Where id = ?"
	stmt, err := MysqlDB.Prepare(sql)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(empId)
	if err != nil {
		return err
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func GetEmployee(empId string) (*Employee, error) {
	var name string
	var id string
	var salary string
	var age string

	err := MysqlDB.QueryRow("SELECT id, name, age, salary FROM employees WHERE id = ?", empId).
		Scan(&id, &name, &salary, &age)
	if err != nil {
		return nil, err
	}
	return &Employee{Id: id, Name: name, Salary: salary, Age: age}, nil
}

func UpdateEmployee(empId string, emp *Employee) error {
	sqlStatement := "UPDATE employees SET name=?, salary=?, age=? WHERE id=? "
	res, err := MysqlDB.Query(sqlStatement, emp.Name, emp.Salary, emp.Age, empId)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func ListEmployee() ([]Employee, error) {
	sqlStatement := "SELECT id, name, salary, age FROM employees order by id"
	rows, err := MysqlDB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Employee{}

	for rows.Next() {
		employee := Employee{}
		err := rows.Scan(&employee.Id, &employee.Name, &employee.Salary, &employee.Age)
		if err != nil {
			return nil, err
		}
		result = append(result, employee)
	}
	return result, nil
}
