package constant

type Category string

const Storage Category = "storage"
const General Category = "general"
const Network Category = "network"
const KVStore Category = "kvstore"
const Trigger Category = "trigger"
const Service Category = "service"
const SQLStore Category = "sqlstore"
const Database Category = "database"
const CacheStore Category = "cachestore"

// CategoryString returns category name
func CategoryString(category Category) string {
	return string(category)
}
