package constant

type Category string

const CategoryStorage Category = "storage"
const CategoryGeneral Category = "general"
const CategoryNetwork Category = "network"
const CategoryKVStore Category = "kvstore"
const CategoryTrigger Category = "trigger"
const CategoryService Category = "service"
const CategorySQLStore Category = "sqlstore"
const CategoryDatabase Category = "database"
const CategoryConsumer Category = "consumer"
const CategoryCacheStore Category = "cachestore"
const CategoryAuthorizer Category = "authorizer"
const CategoryVM Category = "vm"
const CategoryLibrary Category = "library"

// CategoryString returns category name
func CategoryString(category Category) string {
	return string(category)
}
