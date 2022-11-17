package result

/* TotalResult Model used to obtain the total records from a query */
type TotalResult struct {
	Total int64 `bson:"total"`
}
