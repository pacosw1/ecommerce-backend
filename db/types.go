package db

//Product product schema
type Product struct {
	Name        string   `json:"name"`        //required
	Stock       uint16   `json:"stock"`       //required
	Description string   `json:"description"` //required
	Price       float32  `json:"price"`       //required
	SalePrice   float32  `json:"salePrice,omitempty"`
	Images      []string `json:"images"` //required
	Tags        []string `json:"tags,omitempty"`
	Specs       *Specs   `json:"specs,omitempty"`
}

//Specs serves to calculate shipping price and package dimensions
type Specs struct {
	Weigth float32
	Length float32
	Width  float32
	Height float32
	Unit   string
}

//User serves to store user data
type User struct {
	Email        string
	Username     string
	PasswordHash string
}

//Category serves to organize products
type Category struct {
	Name        string
	Description string
	ParentID    string
}

//Order store purchase and shipping info
type Order struct {
	Date     uint64
	Subtotal float32
	Total    float32
	items    []*Product
	Status   string
	Shipping *Shipping
	UserID   string
}

//Shipping serves to store shipping data like address and tracking number
type Shipping struct {
	Name            string
	Address         *Address
	TrackingNumber  string
	ShippingCompany string
}

//Address serves to store user shipping address
type Address struct {
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Col        string
	Country    string
}

//DiscountCode code with discount
type DiscountCode struct {
	Name       string
	Code       string
	PercentOff float32
	MinPrice   float32
	Uses       uint32
	Active     bool
}
