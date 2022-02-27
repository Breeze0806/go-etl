package plugin

import (
	_ "github.com/Breeze0806/go-etl/datax/plugin/reader/csv"
	_ "github.com/Breeze0806/go-etl/datax/plugin/reader/mysql"
	_ "github.com/Breeze0806/go-etl/datax/plugin/reader/postgres"
	_ "github.com/Breeze0806/go-etl/datax/plugin/reader/xlsx"

	_ "github.com/Breeze0806/go-etl/datax/plugin/writer/mysql"
	_ "github.com/Breeze0806/go-etl/datax/plugin/writer/postgres"
)
