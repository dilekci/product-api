package infrastructure

func setupFullTestData() {
	TruncateTestData(ctx, dbPool)
	InsertTestCategories(ctx, dbPool)
	InsertTestProducts(ctx, dbPool)
}

func setupProductsOnly() {
	TruncateTestData(ctx, dbPool)
	InsertTestCategories(ctx, dbPool)
	InsertTestProducts(ctx, dbPool)
}

func clearTestData() {
	TruncateTestData(ctx, dbPool)
}
