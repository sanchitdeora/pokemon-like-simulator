package data

type Item struct {
	Name        string       `json:"name"`
	Count       int          `json:"count"`
	Category    ItemCategory `json:"category"`
	Cost        int          `json:"cost"`
	Attributes  int          `json:"attributes"`
	Description string       `json:"description"`
}

type ItemCategory string

const (
	MedicalItems ItemCategory = "MedicalItems"
	PokeBalls    ItemCategory = "PokeBalls"
)

type BadgeType struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

func BagContainsItem(bag []*Item, item *Item) bool {
	return GetItemFromBag(bag, item) != nil
}

func GetItemFromBag(bag []*Item, item *Item) *Item {
	for _, i := range bag {
		if i.Name == item.Name {
			return i
		}
	}
	return nil
}
