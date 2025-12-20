package repositories

import "github.com/Masterminds/squirrel"

// sq specifies the placeholder format for SQL queries using dollar signs.
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
