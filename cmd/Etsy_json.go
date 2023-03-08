package cmd

type Etsy_Response struct {
	Count   int           `json:"count"`
	Results []Etsy_Result `json:"results"`
}

type Etsy_Result struct {
	ListingID                   int        `json:"listing_id"`
	UserID                      int        `json:"user_id"`
	ShopID                      int        `json:"shop_id"`
	Title                       string     `json:"title"`
	Description                 string     `json:"description"`
	State                       string     `json:"state"`
	CreationTimestamp           int        `json:"creation_timestamp"`
	CreatedTimestamp            int        `json:"created_timestamp"`
	EndingTimestamp             int        `json:"ending_timestamp"`
	OriginalCreationTimestamp   int        `json:"original_creation_timestamp"`
	LastModifiedTimestamp       int        `json:"last_modified_timestamp"`
	UpdatedTimestamp            int        `json:"updated_timestamp"`
	StateTimestamp              int        `json:"state_timestamp"`
	Quantity                    int        `json:"quantity"`
	ShopSectionID               int        `json:"shop_section_id"`
	FeaturedRank                int        `json:"featured_rank"`
	URL                         string     `json:"url"`
	NumFavorers                 int        `json:"num_favorers"`
	NonTaxable                  bool       `json:"non_taxable"`
	IsTaxable                   bool       `json:"is_taxable"`
	IsCustomizable              bool       `json:"is_customizable"`
	IsPersonalizable            bool       `json:"is_personalizable"`
	PersonalizationIsRequired   bool       `json:"personalization_is_required"`
	PersonalizationCharCountMax int        `json:"personalization_char_count_max"`
	PersonalizationInstructions string     `json:"personalization_instructions"`
	ListingType                 string     `json:"listing_type"`
	Tags                        []string   `json:"tags"`
	Materials                   []string   `json:"materials"`
	ShippingProfileID           int        `json:"shipping_profile_id"`
	ReturnPolicyID              int        `json:"return_policy_id"`
	ProcessingMin               int        `json:"processing_min"`
	ProcessingMax               int        `json:"processing_max"`
	WhoMade                     string     `json:"who_made"`
	WhenMade                    string     `json:"when_made"`
	IsSupply                    bool       `json:"is_supply"`
	ItemWeight                  int        `json:"item_weight"`
	ItemWeightUnit              string     `json:"item_weight_unit"`
	ItemLength                  int        `json:"item_length"`
	ItemWidth                   int        `json:"item_width"`
	ItemHeight                  int        `json:"item_height"`
	ItemDimensionsUnit          string     `json:"item_dimensions_unit"`
	IsPrivate                   bool       `json:"is_private"`
	Style                       []string   `json:"style"`
	FileData                    string     `json:"file_data"`
	HasVariations               bool       `json:"has_variations"`
	ShouldAutoRenew             bool       `json:"should_auto_renew"`
	Language                    string     `json:"language"`
	Price                       Etsy_Price `json:"price"`
	TaxonomyID                  int        `json:"taxonomy_id"`
}

type Etsy_Price struct {
	Amount       int    `json:"amount"`
	Divisor      int    `json:"divisor"`
	CurrencyCode string `json:"currency_code"`
}
