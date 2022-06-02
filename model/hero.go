package model

type Hero struct {
	Name      string `csv:"name"`
	Gender    string `csv:"Gender"`
	EyeColor  string `csv:"Eye color"`
	Race      string `csv:"Race"`
	HairColor string `csv:"Hair color"`
	Height    string `csv:"Height"`
	Publisher string `csv:"Publisher"`
	SkinColor string `csv:"Skin color"`
	Alignment string `csv:"Alignment"`
	Weight    string `csv:"Weight"`
}
