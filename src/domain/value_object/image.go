package value_object

type Image string

func ParseImage(uri string) (Image, error) {
	return Image(uri), nil
}

func (i Image) String() string {
	return string(i)
}
