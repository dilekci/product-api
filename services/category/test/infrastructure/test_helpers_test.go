package infrastructure

func setupFullTestData() {
	TruncateTestData(ctx, dbPool)
	InsertTestCategories(ctx, dbPool)
}

func clearTestData() {
	TruncateTestData(ctx, dbPool)
}
