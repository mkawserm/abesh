package constant

type Category string

const KVStore Category = "kvstore"
const Trigger Category = "trigger"
const Service Category = "service"
const SQLStore Category = "sqlstore"
const CacheStore Category = "cachestore"

// CategoryString returns category name
func CategoryString(category Category) string {
	return string(category)
}
