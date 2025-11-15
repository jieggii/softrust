package ports

type GenerateProductReport struct {
	// Query can be anything related to the product:
	// e.g. product title or link to product's website.
	Query string
}

type GenerateVendorReport struct {
	// Query can be anything related to the vendor:
	// e.g. vendor name or link to vendor's website.
	Query string
}
