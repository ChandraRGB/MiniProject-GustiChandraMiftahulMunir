package domain

import "time"

// User represents the user table.
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Nama         string    `gorm:"column:nama;size:255;not null"`
	KataSandi    string    `gorm:"column:kata_sandi;size:255;not null"`
	NoTelp       string    `gorm:"column:notelp;size:255;unique;not null"`
	TanggalLahir time.Time `gorm:"column:tanggal_lahir;type:date"`
	JenisKelamin string    `gorm:"column:jenis_kelamin;size:255"`
	Tentang      string    `gorm:"column:tentang;type:text"`
	Pekerjaan    string    `gorm:"column:pekerjaan;size:255"`
	Email        string    `gorm:"column:email;size:255;unique;not null"`
	IDProvinsi   string    `gorm:"column:id_provinsi;size:255"`
	IDKota       string    `gorm:"column:id_kota;size:255"`
	IsAdmin      bool      `gorm:"column:is_admin"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relations
	Toko      *Toko   `gorm:"foreignKey:UserID"`
	Alamat    []Alamat `gorm:"foreignKey:UserID"`
	Transaksi []Trx    `gorm:"foreignKey:UserID"`
}

func (User) TableName() string { return "user" }

// Toko represents the toko table.
type Toko struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"column:id_user;not null"`
	NamaToko  string    `gorm:"column:nama_toko;size:255;not null"`
	UrlFoto   string    `gorm:"column:url_foto;size:255"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	User   *User    `gorm:"foreignKey:UserID;references:ID"`
	Produk []Produk `gorm:"foreignKey:TokoID"`
}

func (Toko) TableName() string { return "toko" }

// Alamat represents the alamat table.
type Alamat struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uint      `gorm:"column:id_user;not null"`
	JudulAlamat  string    `gorm:"column:judul_alamat;size:255;not null"`
	NamaPenerima string    `gorm:"column:nama_penerima;size:255;not null"`
	NoTelp       string    `gorm:"column:no_telp;size:255;not null"`
	DetailAlamat string    `gorm:"column:detail_alamat;size:255;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (Alamat) TableName() string { return "alamat" }

// Category represents the category table.
type Category struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Nama      string    `gorm:"column:nama_category;size:255;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Produk []Produk `gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string { return "category" }

// Produk represents the produk table.
type Produk struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	NamaProduk    string    `gorm:"column:nama_produk;size:255;not null"`
	Slug          string    `gorm:"column:slug;size:255;not null"`
	HargaReseller string    `gorm:"column:harga_reseller;size:255;not null"`
	HargaKonsumen string    `gorm:"column:harga_konsumen;size:255;not null"`
	Stok          int       `gorm:"column:stok;not null"`
	Deskripsi     string    `gorm:"column:deskripsi;type:text"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
	TokoID        uint      `gorm:"column:id_toko;not null"`
	CategoryID    uint      `gorm:"column:id_category;not null"`

	Toko       Toko        `gorm:"foreignKey:TokoID;references:ID"`
	Category   Category    `gorm:"foreignKey:CategoryID;references:ID"`
	FotoProduk []FotoProduk `gorm:"foreignKey:ProdukID"`
	LogProduk  []LogProduk  `gorm:"foreignKey:ProdukID"`
}

func (Produk) TableName() string { return "produk" }

// FotoProduk represents the foto_produk table.
type FotoProduk struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	ProdukID  uint      `gorm:"column:id_produk;not null"`
	URL       string    `gorm:"column:url;size:255;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Produk Produk `gorm:"foreignKey:ProdukID;references:ID"`
}

func (FotoProduk) TableName() string { return "foto_produk" }

// Trx represents the trx table.
type Trx struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	UserID             uint      `gorm:"column:id_user;not null"`
	AlamatPengirimanID uint      `gorm:"column:alamat_pengiriman;not null"`
	HargaTotal         int       `gorm:"column:harga_total;not null"`
	KodeInvoice        string    `gorm:"column:kode_invoice;size:255;not null"`
	MethodBayar        string    `gorm:"column:method_bayar;size:255;not null"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`

	User      User        `gorm:"foreignKey:UserID;references:ID"`
	Alamat    Alamat      `gorm:"foreignKey:AlamatPengirimanID;references:ID"`
	DetailTrx []DetailTrx `gorm:"foreignKey:TrxID"`
}

func (Trx) TableName() string { return "trx" }

// LogProduk represents the log_produk table.
type LogProduk struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	ProdukID      uint      `gorm:"column:id_produk;not null"`
	NamaProduk    string    `gorm:"column:nama_produk;size:255;not null"`
	Slug          string    `gorm:"column:slug;size:255;not null"`
	HargaReseller string    `gorm:"column:harga_reseller;size:255;not null"`
	HargaKonsumen string    `gorm:"column:harga_konsumen;size:255;not null"`
	Deskripsi     string    `gorm:"column:deskripsi;type:text"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
	TokoID        uint      `gorm:"column:id_toko;not null"`
	CategoryID    uint      `gorm:"column:id_category;not null"`

	Produk    Produk      `gorm:"foreignKey:ProdukID;references:ID"`
	Toko      Toko        `gorm:"foreignKey:TokoID;references:ID"`
	Category  Category    `gorm:"foreignKey:CategoryID;references:ID"`
	DetailTrx []DetailTrx `gorm:"foreignKey:LogProdukID"`
}

func (LogProduk) TableName() string { return "log_produk" }

// DetailTrx represents the detail_trx table.
type DetailTrx struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	TrxID       uint      `gorm:"column:id_trx;not null"`
	LogProdukID uint      `gorm:"column:id_log_produk;not null"`
	TokoID      uint      `gorm:"column:id_toko;not null"`
	Kuantitas   int       `gorm:"column:kuantitas;not null"`
	HargaTotal  int       `gorm:"column:harga_total;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`

	Trx       Trx       `gorm:"foreignKey:TrxID;references:ID"`
	LogProduk LogProduk `gorm:"foreignKey:LogProdukID;references:ID"`
	Toko      Toko      `gorm:"foreignKey:TokoID;references:ID"`
}

func (DetailTrx) TableName() string { return "detail_trx" }
