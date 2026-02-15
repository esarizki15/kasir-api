package models

type SalesSummary struct {
	TotalRevenue   int         `json:"total_revenue"`
	TotalTransaksi int         `json:"total_transaksi"`
	ProdukTerlaris *BestSeller `json:"produk_terlaris,omitempty"`
}

type BestSeller struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
