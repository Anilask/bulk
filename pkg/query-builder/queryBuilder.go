package querybuilder

import (
	"fmt"
	"strconv"
	"strings"
)

type QueryBuilder struct {
	parametersList      []string
	baseCommandTemplate string
	query               string
	currentParamIndex   int
}

func NewQueryBuilder(baseCommandTemplate string) *QueryBuilder {
	queryBuilder := QueryBuilder{}
	queryBuilder.parametersList = []string{}
	queryBuilder.query = baseCommandTemplate + " "
	queryBuilder.currentParamIndex = 1
	return &queryBuilder
}

type JoinFilter struct {
	joinTableName  string
	queryParamName string
	joinColumnName string
}

// NewJoinFilter creates a new JoinFilter which maps a specific filterCriteria using queryParamName to perform a join
// e.g: queryBuilder.GenerateFilterQuery(filterCriteria, []utils.JoinFilter{utils.NewJoinFilter("role", "user_role", "user_id")}).
// means that if a filterCriteria with the name "role" is present, the query will append a join sub-query on the table "user_role" using the column "user_id"
// user_id IN (Select user_id from user_role)
func NewJoinFilter(queryParamName, joinTableName, joinColumn string) JoinFilter {
	joinFilter := JoinFilter{
		joinTableName:  joinTableName,
		queryParamName: queryParamName,
		joinColumnName: joinColumn,
	}
	return joinFilter
}

// GetQuery returns the query string
func (queryBuilder *QueryBuilder) GetQuery() string {
	return queryBuilder.query
}
func (queryBuilder *QueryBuilder) GetParametersList() []string {
	return queryBuilder.parametersList
}

// BindParametersList expects a query and a parametersList as arguments
// The method validates if the given arguments can be bound
// returns the new query after replacing the query's placeholders with the parameters or an error otherwise
func (queryBuilder *QueryBuilder) BindParametersList(query string, parameters []string) (string, error) {
	placeholderCount := strings.Count(query, "$")
	if placeholderCount != len(parameters) {
		return "", fmt.Errorf("failed to bind parameters to query, expected %d parameters, found %d", placeholderCount, len(parameters))
	}
	for index, parameter := range parameters {
		query = strings.Replace(query, "$"+strconv.Itoa(index+1), "'"+parameter+"'", 1)
	}
	return query, nil
}

// BindParameter expects a query and a parameterValue and placeholderIndex as arguments
// The method validates if the given placeholder index is found,
// returns the new query after replacing the query's placeholder at the given index with parameterValue or an error otherwise
func (queryBuilder *QueryBuilder) BindParameter(query string, placeholderIndex int, parameterValue string) (string, error) {
	currentIndex := 0
	errorMessage := ""
	if placeholderIndex < 1 {
		return query, fmt.Errorf("parameter index should be > 0")
	}
	placeholderCount := strings.Count(query, "$")
	if placeholderCount < placeholderIndex {
		return query, fmt.Errorf("placeholderIndex is out of bounds: expected " + strconv.Itoa(placeholderCount) + " found " + strconv.Itoa(placeholderIndex))
	}
	for parameterIndex := 1; parameterIndex <= placeholderIndex; parameterIndex++ {
		characterIndex := strings.Index(query[currentIndex:], "$"+strconv.Itoa(placeholderIndex))
		if characterIndex < 0 {
			errorMessage = "query statement has no placeholder for the provided parameter"
			break
		}
		currentIndex += characterIndex
		if placeholderIndex == parameterIndex {
			return query[:currentIndex] + parameterValue + query[currentIndex+len("$"+strconv.Itoa(placeholderIndex)):], nil
		}
		currentIndex += len("$" + strconv.Itoa(placeholderIndex))
	}
	return query, fmt.Errorf("failed to bind parameter: %s", errorMessage)
}

// Build validates if the given placeholder index is found,
// returns the new query after replacing the query's placeholders with the parameters or an error otherwise
func (queryBuilder *QueryBuilder) Build() error {
	placeholderCount := strings.Count(queryBuilder.query, "$")
	if placeholderCount != len(queryBuilder.parametersList) {
		return fmt.Errorf("failed to bind parameters to query, expected %d parameters, found %d", placeholderCount, len(queryBuilder.parametersList))
	}
	for index, parameter := range queryBuilder.parametersList {
		queryBuilder.query = strings.Replace(queryBuilder.query, "$"+strconv.Itoa(index+1), "'"+parameter+"'", 1)
	}
	return nil
}

// // Select appends the Select clause along with the required columns
// // the method expect the tableName string and column Names []string
// func (queryBuilder *QueryBuilder) Select(tableName string, columnNames []string) *QueryBuilder {
// 	queryBuilder.query += "SELECT "
// 	for index, _ := range columnNames {
// 		if index == 0 {
// 			queryBuilder.query += columnNames[index] + " "
// 		} else {
// 			queryBuilder.query += ", " + columnNames[index] + " "
// 		}
// 	}
// 	queryBuilder.query += "FROM " + tableName + " "
// 	return queryBuilder
// }

// // Create appends the Create clause along with the required columns
// // the method expect the tableName string, columnNames and columnValues []string
// func (queryBuilder *QueryBuilder) Create(tableName string, columnNames, columnValues []string) *QueryBuilder {
// 	queryBuilder.query = "INSERT INTO " + tableName + " "
// 	columnsString := "("
// 	valuesString := "VALUES ("
// 	for index, _ := range columnNames {
// 		if index == 0 {
// 			columnsString += columnNames[index]
// 			valuesString += "'" + columnValues[index] + "'"
// 		} else {
// 			columnsString += "," + columnNames[index]
// 			valuesString += ",'" + columnValues[index] + "'"
// 		}
// 	}
// 	queryBuilder.query += columnsString + ") " + valuesString + " ); "
// 	return queryBuilder
// }

// // Update appends the Update clause along with the required columnNames, columnValues to set
// // the method expect the tableName string and column Names []string
// func (queryBuilder *QueryBuilder) Update(tableName string, columnNames, columnValues []string) *QueryBuilder {
// 	queryBuilder.query += "UPDATE " + tableName + " SET "
// 	for index, _ := range columnNames {
// 		if index == 0 {
// 			queryBuilder.query += columnNames[index] + " =$" + strconv.Itoa(queryBuilder.currentParamIndex) + " "
// 		} else {
// 			queryBuilder.query += ", " + columnNames[index] + " =$" + strconv.Itoa(queryBuilder.currentParamIndex) + " "
// 		}
// 		queryBuilder.currentParamIndex++
// 	}
// 	queryBuilder.parametersList = append(queryBuilder.parametersList, columnValues...)
// 	return queryBuilder
// }

// // Delete appends the Delete clause
// // the method expect the tableName string
// func (queryBuilder *QueryBuilder) Delete(tableName string) *QueryBuilder {
// 	queryBuilder.query = "DELETE FROM " + tableName + " "
// 	return queryBuilder
// }

// Where appends the where clause and columnName to the query
// The method expects columnName string
func (queryBuilder *QueryBuilder) Where(columnName string) *QueryBuilder {
	queryBuilder.query += "WHERE " + columnName + " "
	return queryBuilder
}

// // Exists appends the EXISTS clause to the query to the query
// func (queryBuilder *QueryBuilder) Exists() *QueryBuilder {
// 	queryBuilder.query += "EXISTS "
// 	return queryBuilder
// }

// // StartSubQuery opens a sub-query within the base query
// // The method expects the tableName, subQueryCommandTemplate and columnName to use in the subQuery
// func (queryBuilder *QueryBuilder) StartSubQuery() *QueryBuilder {
// 	queryBuilder.query += "(" // + subQueryCommandTemplate + " WHERE " + tableName + "." + columnName
// 	return queryBuilder
// }

// // EndSubQuery closes the opened sub-query within the base query
// func (queryBuilder *QueryBuilder) EndSubQuery() *QueryBuilder {
// 	queryBuilder.query += ") "
// 	return queryBuilder
// }

// // From appends the FROM clause to the query and tableName to the query
// // The method expects the tableName string as an argument
// func (queryBuilder *QueryBuilder) From(tableName string) *QueryBuilder {
// 	queryBuilder.query += "FROM " + tableName + " "
// 	return queryBuilder
// }

// And appends the AND clause and columnName to the query
// The method expects columnName string
func (queryBuilder *QueryBuilder) And(columnName string) *QueryBuilder {
	queryBuilder.query += "AND " + columnName + " "
	return queryBuilder
}
func (queryBuilder *QueryBuilder) Set(columnName string) *QueryBuilder {
	queryBuilder.query += "SET " + columnName + " "
	return queryBuilder
}
func (queryBuilder *QueryBuilder) AddComma(columnName string) *QueryBuilder {
	queryBuilder.query += ", " + columnName + " "
	return queryBuilder
}

// // Or appends the OR clause  and columnName to the query
// // The method expects columnName string
// func (queryBuilder *QueryBuilder) Or(columnName string) *QueryBuilder {
// 	queryBuilder.query += "OR " + columnName + " "
// 	return queryBuilder
// }

// // Limit appends the LIMIT clause to the query and expects the limit integer as an argument
func (queryBuilder *QueryBuilder) Limit(limit int) *QueryBuilder {
	queryBuilder.query += "LIMIT " + strconv.Itoa(limit) + " "
	return queryBuilder
}

// // Offset appends the OFFSET clause to the query and expects the offset integer as an argument
func (queryBuilder *QueryBuilder) Offset(offset int) *QueryBuilder {
	queryBuilder.query += "OFFSET " + strconv.Itoa(offset) + " "
	return queryBuilder
}

// Sort appends the ORDER BY clause to the query and expects columnName string to sort by and mode boolean (true = DESC, false = ASC)
// // if columnName is not provided (empty string) the method returns without appending the Sort clause
func (queryBuilder *QueryBuilder) Sort(columnName string, mode bool) *QueryBuilder {
	if columnName == "" {
		return queryBuilder
	}
	queryBuilder.query += "ORDER BY " + columnName + " "
	if mode {
		queryBuilder.query += "DESC "
	}
	return queryBuilder
}

// Equals appends the Equals clause to the query
// The method expects columnValue string
func (queryBuilder *QueryBuilder) Equals(columnValues string) *QueryBuilder {
	queryBuilder.parametersList = append(queryBuilder.parametersList, columnValues)
	queryBuilder.query += " = $" + strconv.Itoa(queryBuilder.currentParamIndex) + " "
	queryBuilder.currentParamIndex++
	return queryBuilder
}

func (queryBuilder *QueryBuilder) GreaterThanEquals(columnValues string) *QueryBuilder {
	queryBuilder.parametersList = append(queryBuilder.parametersList, columnValues)
	queryBuilder.query += ">= $" + strconv.Itoa(queryBuilder.currentParamIndex) + " "
	queryBuilder.currentParamIndex++
	return queryBuilder
}

func (queryBuilder *QueryBuilder) LessThanEquals(columnValues string) *QueryBuilder {
	queryBuilder.parametersList = append(queryBuilder.parametersList, columnValues)
	queryBuilder.query += "<= $" + strconv.Itoa(queryBuilder.currentParamIndex) + " "
	queryBuilder.currentParamIndex++
	return queryBuilder
}

func (queryBuilder *QueryBuilder) Like(columnValues string) *QueryBuilder {
	queryBuilder.parametersList = append(queryBuilder.parametersList, columnValues)
	queryBuilder.query += " like $" + strconv.Itoa(queryBuilder.currentParamIndex) + " "
	queryBuilder.currentParamIndex++
	return queryBuilder
}

func (queryBuilder *QueryBuilder) ReverseIn(columnValues ...string) *QueryBuilder {

	var columnNames string
	for _, val := range columnValues {
		columnNames += val + ","
	}
	columnNames = columnNames[:len(columnNames)-1]
	queryBuilder.query += " in (" + columnNames + ") "
	return queryBuilder
}

// //GenerateSelectStatement generates a select statement
// //The method expects selectedColumnName,tableName as strings, columnNames,columnValues as string arrays
// func (queryBuilder *QueryBuilder) GenerateSelectStatement(selectedColumnName, tableName string, columnNames, columnValues []string) *QueryBuilder {
// 	queryBuilder.Select(tableName, []string{selectedColumnName})
// 	for index, _ := range columnNames {
// 		if index == 0 {
// 			queryBuilder.Where(columnNames[index]).Equals(columnValues[index])
// 		} else {
// 			queryBuilder.And(columnNames[index]).Equals(columnValues[index])
// 		}
// 	}
// 	return queryBuilder
// }

// func (queryBuilder *QueryBuilder) Between(start string, end string) *QueryBuilder {
// 	queryBuilder.query += " BETWEEN '" + start + "' AND '" + end + "' "
// 	return queryBuilder
// }
