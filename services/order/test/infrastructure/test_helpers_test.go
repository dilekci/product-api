package infrastructure

func setupOrdersOnly() {
	TruncateTestData(ctx, dbPool)
	InsertTestOrders(ctx, dbPool)
}

func clearTestData() {
	TruncateTestData(ctx, dbPool)
}
